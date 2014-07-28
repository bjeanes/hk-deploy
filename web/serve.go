package main

import (
	"fmt"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func curlHandler(w http.ResponseWriter, r *http.Request) {
	policy := NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, policy.ToCurl())
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	policy := NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, policy.ToJsonResponse())
}

func Serve() {
	if port == "" {
		port = "5000"
	}

	listen := fmt.Sprintf(":%s", port)
	fmt.Println("Listening on " + listen)

	http.HandleFunc("/curl", curlHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(listen, nil)
}
