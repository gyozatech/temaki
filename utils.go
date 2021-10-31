package temaki

import (
	"net/http"
)

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func GetPathParam(r *http.Request, param string) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	pathParamsMap := r.Context().Value("pathParamsMap").(map[string]int)
	return fields[pathParamsMap[param]]
}
