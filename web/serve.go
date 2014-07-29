package main

import (
	"fmt"
	"net/http"
	"os"
)

var port = os.Getenv("PORT")

func formCurlHandler(w http.ResponseWriter, r *http.Request) {
	policy := NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, policy.ToCurl())
}

func formJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	policy := NewPolicy()
	fmt.Printf("Serving upload policy for /%s\n", policy.Key())
	fmt.Fprintln(w, policy.ToJsonResponse())
}
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/bjeanes/hk-deploy", http.StatusSeeOther)
}

func Serve() {
	if port == "" {
		port = "5000"
	}

	listen := fmt.Sprintf(":%s", port)
	fmt.Println("Listening on " + listen)

	// original routes:
	http.HandleFunc("/curl", formCurlHandler)
	http.HandleFunc("/slot", formJsonHandler)

	// new routes:
	http.HandleFunc("/form.curl", formCurlHandler)
	http.HandleFunc("/form.json", formJsonHandler)

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(listen, nil)
}
