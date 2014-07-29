package s3

import (
	"net/http"
	. "testing"
	"time"
)

// Based on: http://goo.gl/dTg2U1
func TestRequestSigning(t *T) {
	region := "us-east-1"
	keyId := "AKIAIOSFODNN7EXAMPLE"
	secretKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	bucketName := "examplebucket"
	expireIn := uint(86400)
	time, _ := time.Parse(time.RFC1123, "Fri, 24 May 2013 00:00:00 GMT")
	origReq, _ := http.NewRequest("GET", "https://"+bucketName+".s3.amazonaws.com/test.txt", nil)

	bucket := Bucket{S3{Auth{keyId, secretKey}, region}, bucketName}

	newReq, ok := bucket.SignRequest(origReq, time, expireIn)
	if !ok {
		t.Error("Request signing did not return OK")
	}

	expectedUrl := "https://examplebucket.s3.amazonaws.com/test.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20130524%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20130524T000000Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=aeeed9bbccd4d02ee5c0109b86d86835f995330da4c265957d157751f604d404"

	if actualUrl := newReq.URL.String(); actualUrl != expectedUrl {
		t.Errorf("Incorrect URL generated:\n  expected: %s\n  actual:   %s", expectedUrl, actualUrl)
	}
}
