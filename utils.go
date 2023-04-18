package temaki

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// GetField fetched a field from the request context
func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

// GetPathParam fetched the specified path param from the request
func GetPathParam(r *http.Request, param string) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	pathParamsMap := r.Context().Value("pathParamsMap").(map[string]int)
	return fields[pathParamsMap[param]]
}

// GetBasicToken fetches a basic authorization token from the http request
func GetBasicToken(r *http.Request) (username, password string, err error) {
	if r == nil {
		return "", "", fmt.Errorf("invalid HTTP request")
	}
	value := r.Header.Get("Authorization")
	if strings.HasPrefix(value, "Basic ") || strings.HasPrefix(value, "basic ") {
		split := strings.Split(value, " ")
		if len(split) == 2 {
			decoded, err := base64.StdEncoding.DecodeString(split[1])
			if err != nil {
				return "", "", fmt.Errorf("basic Authorization token is malformed")
			}
			split = strings.Split(string(decoded), ":")
			if len(split) == 2 {
				return split[0], split[1], nil
			}
		}
	}
	return "", "", fmt.Errorf("basic Authorization token is missing")
}

// GetBearerToken fetches a bearer authorization token from the http request
func GetBearerToken(r *http.Request) (string, error) {
	if r == nil {
		return "", fmt.Errorf("invalid HTTP request")
	}
	value := r.Header.Get("Authorization")
	if strings.HasPrefix(value, "Bearer ") || strings.HasPrefix(value, "bearer ") {
		split := strings.Split(value, " ")
		if len(split) == 2 {
			return split[1], nil
		}
	}
	return "", fmt.Errorf("bearer Authorization token is missing")
}
