package main

import (
	// s3 "github.com/mitchellh/goamz/s3"
	"code.google.com/p/go-uuid/uuid"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

var (
	bucket = os.Getenv("AWS_S3_BUCKET")
	keyId  = os.Getenv("AWS_ACCESS_KEY_ID")
	secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	action = "https://" + bucket + ".s3.amazonaws.com/"
	expire = 10 * time.Minute
)

const (
	region      = "us-east-1" // FIXME: can be derived from bucket?
	acl         = "private"
	algorithm   = "AWS4-HMAC-SHA256"
	contentType = "application/x-compressed"

	// time layouts based on reference time (see pkg "time")
	iso8601 = "20060102T150405Z0700"
	short   = "20060102"
)

func sign(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func creds(now time.Time) (credential string, signingKey []byte) {
	credential = keyId + "/" + now.Format(short) + "/" + region + "/s3/aws4_request"

	// http://goo.gl/cPOyvG
	signingKey = []byte("AWS4" + secret)
	signingKey = sign(signingKey, []byte(now.Format(short)))
	signingKey = sign(signingKey, []byte(region))
	signingKey = sign(signingKey, []byte("s3"))
	signingKey = sign(signingKey, []byte("aws4_request"))

	return
}

func main() {
	now := time.Now().UTC()
	isoNow := now.Format(iso8601)
	expiry := now.Add(expire).Format(time.RFC3339)
	key := now.Format(short) + "/" + uuid.New() + ".tgz"

	credential, signingKey := creds(now)

	policy := `{
	"expiration": "` + expiry + `",
	"conditions": [
		{"acl": "` + acl + `" },
		{"bucket": "` + bucket + `" },
		{"key": "` + key + `"},
		{"x-amz-date": "` + isoNow + `"},
		{"x-amz-credential": "` + credential + `"},
		{"x-amz-algorithm": "` + algorithm + `"},
		{"content-type": "` + contentType + `"}
	]
}`

	encodedPolicy := base64.StdEncoding.EncodeToString([]byte(policy))
	signature := hex.EncodeToString(sign(signingKey, []byte(encodedPolicy)))

	var response string

	if len(os.Args) < 2 || os.Args[1] != "curl" {
		response = `{
	"action": "` + action + `",
	"fields": {
		"key": "` + key + `",
		"acl": "` + acl + `",
		"Content-Type": "` + contentType + `",
		"X-Amz-Credential": "` + credential + `",
		"X-Amz-Algorithm": "` + algorithm + `",
		"X-Amz-Date": "` + isoNow + `",
		"Policy": "` + encodedPolicy + `",
		"X-Amz-Signature": "` + signature + `"
	}
}`
	} else {
		response = `curl ` + action + ` ` +
			`-Fkey="` + key + `" ` +
			`-Facl="` + acl + `" ` +
			`-FContent-Type="` + contentType + `" ` +
			`-FX-Amz-Credential="` + credential + `" ` +
			`-FX-Amz-Algorithm="` + algorithm + `" ` +
			`-FX-Amz-Date="` + isoNow + `" ` +
			`-FPolicy="` + encodedPolicy + `" ` +
			`-FX-Amz-Signature="` + signature + `" ` +
			`-Ffile=@foo.tgz`

	}

	fmt.Println(response)
}
