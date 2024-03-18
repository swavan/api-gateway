package handler

import (
	"encoding/json"
	"net/http"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	payload := new(struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	})
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := a.api.User().
		FindByUsername(
			r.Context(),
			payload.Username)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return
	}

	// Match password and generate token

}
