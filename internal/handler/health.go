package handler

import "net/http"

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("API Gateway is running"))
}
