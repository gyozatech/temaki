package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gyozatech/temaki/reverseproxy"
)

func main() {
	// you must start main_s1.go and main_s2.go first
	os.Setenv("S1_SERVICE_PROXY_URL", "(http://localhost:8081)/s1-service/api/v1/")
	os.Setenv("S2_SERVICE_PROXY_URL", "(http://localhost:8082)/s2-service/api/v1/")

	log.Println("Variable: ", os.Getenv("LOGIN_SERVICE_PROXY_URL"))

	log.Fatal(reverseproxy.NewReverseProxy().
		UseMiddleware(JWTMiddleware).
		UseRequestModifier(modifyRequest).
		UseResponseModifier(modifyResponse).
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

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func modifyResponse(resp *http.Response) error {
	resp.Header.Set("X-Proxy", "Magical")
	return nil
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}
