package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gyozatech/temaki/reverseproxy"
	"github.com/gyozatech/temaki/middlewares"
)

func main() {

	// you must start main_s1.go and main_s2.go first
	os.Setenv("PROXY_RULE_S1", "/s1/>http://localhost:8081")
	os.Setenv("PROXY_RULE_S2", "/s2/>http://localhost:8082")

	routes := reverseproxy.CollectPathPrefixRoutesFromEnvVar()
	// or alternatively, you can initialize directly the PathPrefixRoutesMap: 
	/* 
           routes := reverseproxy.PathPrefixRoutesMap{
		"s1": "http://localhost:8081",
		"s2": "http://localhost:8082",
	}
        */

	log.Fatalf("Server error: %s", reverseproxy.New(routes).
		WithMiddlewares(middlewares.RequestLoggerMiddleware, JWTMiddleware).
		Start(8080))
}

func JWTMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer abcd" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}
