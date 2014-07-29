package s3

import "os"

const (
	region    = "us-east-1"
	algorithm = "AWS4-HMAC-SHA256"

	// time layouts based on reference time (see pkg "time")
	iso8601 = "20060102T150405Z0700"
	short   = "20060102"
)

type Auth struct {
	AccessKey, SecretKey string
}

var EnvAuth = Auth{os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")}

type S3 struct {
	Auth
	Region string
}

type Bucket struct {
	S3
	Name string
}
