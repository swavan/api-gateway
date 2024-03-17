package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(health))
	http.ListenAndServe(":9000", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Alert service is available"))
}
