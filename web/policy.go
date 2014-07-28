package main

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"time"
)

var (
	bucket      = os.Getenv("AWS_S3_BUCKET")
	keyId       = os.Getenv("AWS_ACCESS_KEY_ID")
	secret      = os.Getenv("AWS_SECRET_ACCESS_KEY")
	action      = "https://" + bucket + ".s3.amazonaws.com/"
	expireAfter = 10 * time.Minute
)

func sign(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

type policy struct {
	time time.Time
	uuid string
}

func NewPolicy() policy {
	return policy{time.Now().UTC(), uuid.New()}
}

func (policy policy) expiry() time.Time {
	return policy.time.Add(expireAfter)
}

func (policy policy) Key() string {
	return "uploads/" + policy.time.Format(short) + "/" + policy.uuid + ".tgz"
}

func (policy policy) Credential() string {
	// TODO: don't get keyId from global but add it to policy
	return keyId + "/" +
		policy.time.Format(short) + "/" +
		region + "/s3/aws4_request"
}

func (policy policy) ToJson() string {
	return `{
	"expiration": "` + policy.expiry().Format(time.RFC3339) + `",
	"conditions": [
		{"acl": "` + acl + `" },
		{"bucket": "` + bucket + `" },
		{"key": "` + policy.Key() + `"},
		{"x-amz-date": "` + policy.time.Format(iso8601) + `"},
		{"x-amz-credential": "` + policy.Credential() + `"},
		{"x-amz-algorithm": "` + algorithm + `"},
		{"content-type": "` + contentType + `"}
	]
}`
}

func (policy policy) Encoded() string {
	return base64.StdEncoding.EncodeToString([]byte(policy.ToJson()))
}

func (policy policy) Signature() string {
	return hex.EncodeToString(sign(policy.SigningKey(), []byte(policy.Encoded())))
}

func (policy policy) SigningKey() []byte {
	// http://goo.gl/cPOyvG
	signingKey := []byte("AWS4" + secret)
	signingKey = sign(signingKey, []byte(policy.time.Format(short)))
	signingKey = sign(signingKey, []byte(region))
	signingKey = sign(signingKey, []byte("s3"))
	signingKey = sign(signingKey, []byte("aws4_request"))

	return signingKey
}

func (policy policy) ToJsonResponse() string {
	return `{
	"action": "` + action + `",
	"fields": {
		"key": "` + policy.Key() + `",
		"acl": "` + acl + `",
		"Content-Type": "` + contentType + `",
		"X-Amz-Credential": "` + policy.Credential() + `",
		"X-Amz-Algorithm": "` + algorithm + `",
		"X-Amz-Date": "` + policy.time.Format(iso8601) + `",
		"Policy": "` + policy.Encoded() + `",
		"X-Amz-Signature": "` + policy.Signature() + `"
	}
}`
}

func (policy policy) ToCurl() string {
	return `curl -i "` + action + `" ` +
		`-Fkey="` + policy.Key() + `" ` +
		`-Facl="` + acl + `" ` +
		`-FContent-Type="` + contentType + `" ` +
		`-FX-Amz-Credential="` + policy.Credential() + `" ` +
		`-FX-Amz-Algorithm="` + algorithm + `" ` +
		`-FX-Amz-Date="` + policy.time.Format(iso8601) + `" ` +
		`-FPolicy="` + policy.Encoded() + `" ` +
		`-FX-Amz-Signature="` + policy.Signature() + `" ` +
		`-Ffile=@$FILE_TO_UPLOAD`
}
