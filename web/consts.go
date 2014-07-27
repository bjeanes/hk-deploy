package main

const (
	region      = "us-east-1" // FIXME: can be derived from bucket?
	acl         = "private"
	algorithm   = "AWS4-HMAC-SHA256"
	contentType = "application/x-compressed"

	// time layouts based on reference time (see pkg "time")
	iso8601 = "20060102T150405Z0700"
	short   = "20060102"
)
