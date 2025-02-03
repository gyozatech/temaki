/*
	service S2 must be started before main.go: the reverse proxy is going to proxy the requests prefixed with "s2" to this service 
*/
package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK from s2-service"))
		return
	})

	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from s2-service"))
		return
	})
	
	log.Println("S2 HTTP server started on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
	
}
