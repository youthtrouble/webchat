package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		//not authticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return 
	}
	//success, call the next handler
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r) 
}

//MustAuth is a middleware for authentication
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
			default:
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Auth action %s not supported", action)


	}
}
