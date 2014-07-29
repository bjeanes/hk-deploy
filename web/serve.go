package main

import (
	"fmt"
	"github.com/bjeanes/hk-deploy/policy"
	"github.com/bjeanes/hk-deploy/s3"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func formCurlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/x-shellscript")

	policy := policy.NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, `#!/usr/bin/env sh
FILE_TO_UPLOAD=$1
`+policy.ToCurl())
}

func formJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	policy := policy.NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, policy.ToJsonResponse())
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/bjeanes/hk-deploy", http.StatusSeeOther)
}

var bucket = s3.Bucket{s3.S3{s3.EnvAuth, "us-east-1"}, os.Getenv("AWS_S3_BUCKET")}

func Serve() {
	if port == "" {
		port = "5000"
	}

	listen := fmt.Sprintf(":%s", port)
	fmt.Println("Listening on " + listen)

	// new routes:
	http.HandleFunc("/form.sh", formCurlHandler)
	http.HandleFunc("/form.json", formJsonHandler)
	http.HandleFunc("/slot.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		upload := NewUpload(&bucket)
		fmt.Printf("Serving upload policy for /%s\n", upload.Key())
		fmt.Fprintln(w, upload.ToJson())
	})

	// original routes:
	http.HandleFunc("/curl", formCurlHandler)
	http.HandleFunc("/slot", formJsonHandler)

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(listen, nil)
}
