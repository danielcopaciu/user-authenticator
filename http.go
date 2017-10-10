package main

import (
	"encoding/json"
	"net/http"
)

type httpHandler struct {
	manager
}

func newHTTPHandler(manager manager) httpHandler {
	return httpHandler{manager}
}

func (h httpHandler) doAuthentication(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.authenticate(username, password)

	if err == ErrInvalidUsername || err == ErrInvalidPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"username":  user.username,
		"sessionID": sessionID(),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
