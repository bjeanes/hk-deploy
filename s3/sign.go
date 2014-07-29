package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

import "fmt"

const algo = "AWS4-HMAC-SHA256"

func sign(key, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}

func creds(secretKey, region string, at time.Time) (scope, signingKey string) {
	date := at.Format(short)
	scope = strings.Join([]string{date, region, "s3", "aws4_request"}, "/")
	signingKey = "AWS4" + secretKey
	signingKey = sign(signingKey, at.Format(short))
	signingKey = sign(signingKey, region)
	signingKey = sign(signingKey, "s3")
	signingKey = sign(signingKey, "aws4_request")

	return
}

func escape(s string) string {
	// return strings.NewReplacer("/", "%2F").Replace(s)
	return s
}

func (bucket *Bucket) SignRequest(unsigned *http.Request, at time.Time, expireIn uint) (signed *http.Request, ok bool) {
	signed = new(http.Request)
	*signed = *unsigned // copy unsigned request

	if signed.Method == "" {
		signed.Method = "GET"
	}

	scope, signingKey := creds(bucket.S3.Auth.SecretKey, bucket.Region, at)

	query := unsigned.URL.Query()
	query.Set("X-Amz-Algorithm", algo)
	query.Set("X-Amz-Credential",
		fmt.Sprintf("%s/%s", bucket.S3.Auth.AccessKey, scope))
	query.Set("X-Amz-Date", at.UTC().Format(iso8601))
	query.Set("X-Amz-Expires", fmt.Sprintf("%d", expireIn))

	// TODO: automatically pick multiple headers to sign
	query.Set("X-Amz-SignedHeaders", "host")
	canonicalHeaders := "host:" + signed.URL.Host + "\n"
	signedHeaders := "host"

	canonicalRequest := strings.Join(
		[]string{
			signed.Method,
			escape(signed.URL.Path),
			query.Encode(),
			canonicalHeaders,
			signedHeaders,
			"UNSIGNED-PAYLOAD",
		},
		"\n",
	)

	hashedRequestBytes := sha256.Sum256([]byte(canonicalRequest))

	stringToSign := strings.Join(
		[]string{
			algo,
			at.Format(iso8601),
			scope,
			hex.EncodeToString(hashedRequestBytes[0:32]),
		},
		"\n",
	)

	signature := hex.EncodeToString([]byte(sign(signingKey, stringToSign)))

	// query.Set("X-Amz-Signature", signature)
	signed.URL.RawQuery = query.Encode()
	signed.URL.RawQuery += "&X-Amz-Signature=" + signature // to make it at end

	return signed, true
}
