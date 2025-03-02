/*
	reverse proxy service: it's using env vars to register the routes to the backend services s1 and s2 based on the prefix.
 	after having run:
  		- `go run main_s1.go` on one tab
    		- `go run main_s2.go` on another tab
      	you can run the reverse proxy via:
       		- `go run main.go` onto another tab
	 then you can test the routes with:
  		curl -H "Authorization: Bearer abcd" http://localhost:8080/s1/api/v1/status
    		curl -H "Authorization: Bearer abcd" http://localhost:8080/s1/api/v1/hello
      		curl -H "Authorization: Bearer abcd" http://localhost:8080/s2/api/v1/status
		curl -H "Authorization: Bearer abcd" http://localhost:8080/s2/api/v1/hello
*/
package main

import (
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
