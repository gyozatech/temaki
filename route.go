package temaki

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var rgx = regexp.MustCompile(`\{(.*?)\}`)
var rgxIn = regexp.MustCompile(`\((.*?)\)`)

type Route struct {
	method     string
	regex      *regexp.Regexp
	pathParams map[string]int
	handler    http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) Route {
	pathParamsMap, regexPath := parseURL(pattern)
	if regexPath[0] != '/' {
		regexPath = fmt.Sprintf("/%s", regexPath)
	}
	return Route{method, regexp.MustCompile("^" + regexPath + "$"), pathParamsMap, handler}
}

func parseURL(path string) (map[string]int, string) {
	pathParams := map[string]int{}
	params := rgx.FindAllString(path, -1)

	for i, param := range params {

		oldParam := param
		param = param[1 : len(param)-1]
		pathParamKey := param

		strPattern := "([^/]+)"

		inParam := rgxIn.FindAllString(param, -1)
		if len(inParam) == 1 {
			strPattern = inParam[0]
			pathParamKey = strings.Split(param, strPattern)[0]
			strPattern = inParam[0]
		}
		pathParams[pathParamKey] = i
		path = strings.Replace(path, oldParam, strPattern, 1)
	}

	return pathParams, path
}
