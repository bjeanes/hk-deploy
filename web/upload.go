package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/bjeanes/hk-deploy/s3"
	"net/http"
	"time"
)

const postExpiry = 2 * 60 // 2 minutes
const getExpiry = 20 * 60 // 20 minutes

type upload struct {
	Get, Post string
	key       string
}

func NewUpload(bucket *s3.Bucket) *upload {
	now := time.Now().UTC()

	key := fmt.Sprintf("uploads/%s/%s.tgz",
		now.Format("20060102"),
		uuid.New(),
	)

	var get string
	var post string

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket.Name, key)

	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil
	} else {
		if newReq, ok := bucket.SignRequest(req, now, getExpiry); ok {
			get = newReq.URL.String()
		} else {
			return nil
		}
	}

	if req, err := http.NewRequest("POST", url, nil); err != nil {
		return nil
	} else {
		if newReq, ok := bucket.SignRequest(req, now, postExpiry); ok {
			post = newReq.URL.String()
		} else {
			return nil
		}
	}

	return &upload{get, post, key}
}

func (u upload) Key() string {
	return u.key
}

func (u upload) ToJson() string {
	bytes, _ := json.Marshal(u)
	return string(bytes)
}
