package api

import (
	"bytes"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func useMiddlewares(r *mux.Router) *mux.Router {
	r.Use(loggingMiddleware)
	r.Use(authenticationMiddleware)
	r.Use(headerMiddleware)

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		clientIp := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
		log.Printf(`
IP Address: %s
Requested URI: %s
Method: %s
Body: %s

`,
			clientIp, r.RequestURI, r.Method, body)
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/books") && r.Method == "POST" {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token not found", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func headerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(writer, request)
	})
}
