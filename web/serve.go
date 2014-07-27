package main

import (
	"fmt"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func handler(w http.ResponseWriter, r *http.Request) {
	policy := NewPolicy()
	fmt.Printf("Serving upload policy for %s\n", policy.Key())
	fmt.Fprintf(w, policy.ToJsonResponse())
}

func Serve() {
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/", handler)
	listen := fmt.Sprintf(":%s", port)
	fmt.Println("Listening on " + listen)
	http.ListenAndServe(listen, nil)
}
