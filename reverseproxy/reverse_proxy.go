package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Middleware represent a HTTP middleware for all incoming requests
type Middleware func(next http.Handler) http.Handler

// PathPrefix is the prefix of the path to find a match in the request path
type PathPrefix string

// TargetHost is the destination host to which to proxy the whole request after the PathPrefix removal
type TargetHost string

// PathPrefixRoutesMap is the map of PathPrefixes and TargetHosts to map the available routes for the reverse proxy
type PathPrefixRoutesMap map[PathPrefix]TargetHost

// ReverseProxy is the instance of the reverse proxy server
type ReverseProxy struct {
	middlewares         []Middleware
	pathPrefixRoutesMap PathPrefixRoutesMap
}

// WithMiddlewares allows specifying the http middleware to be applied to all routes
func (r *ReverseProxy) WithMiddlewares(middlewares ...Middleware) *ReverseProxy {
	if r.middlewares == nil {
		r.middlewares = []Middleware{}
	}
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

// CollectPathPrefixRoutesFromEnvVar collects the env vars like:
// PROXY_RULE_WEATHER_API=/weather/>https://api.weather.com
// PROXY_RULE_WEATHER_API=/geo/>https://api.geo.com
// This means that a call to:
//
//	https://<proxy-host>/weather/v1/today?zone=EU
//
// will be sent to:
//
//	https://api.weather.com/v1/today?zone=EU
func CollectPathPrefixRoutesFromEnvVar() PathPrefixRoutesMap {
	var routesMap PathPrefixRoutesMap = make(map[PathPrefix]TargetHost, 0)
	var envVars []string = os.Environ()
	for _, envVar := range envVars {

		envVarName := strings.Split(envVar, "=")[0]
		envVarValue := strings.Split(envVar, "=")[1]

		if strings.HasPrefix(envVarName, "PROXY_RULE_") {
			proxyRule := strings.Split(envVarValue, ">")
			routesMap[PathPrefix(proxyRule[0])] = TargetHost(proxyRule[1])
		}
	}
	return routesMap
}

func (r *ReverseProxy) applyMiddlewares(handler http.Handler) http.Handler {
	for _, m := range r.middlewares {
		handler = m(handler)
	}
	return handler
}

// New is used to create a new instance of a reverse proxy
func New(routes PathPrefixRoutesMap) *ReverseProxy {
	return &ReverseProxy{
		middlewares:         []Middleware{},
		pathPrefixRoutesMap: adaptRoutesMap(routes),
	}
}

func adaptRoutesMap(routes PathPrefixRoutesMap) PathPrefixRoutesMap {
	adaptedRoutesMap := make(PathPrefixRoutesMap, 0)
	for prefix, host := range routes {
		prefixStr := strings.ReplaceAll(string(prefix), " ", "")
		hostStr := strings.ReplaceAll(string(host), " ", "")
		if !strings.HasPrefix(prefixStr, "/") {
			prefixStr = "/" + prefixStr
		}
		if !strings.HasSuffix(prefixStr, "/") {
			prefixStr = prefixStr + "/"
		}
		if !strings.HasPrefix(hostStr, "http") {
			hostStr = "http://" + hostStr
		}
		adaptedRoutesMap[PathPrefix(prefixStr)] = TargetHost(strings.TrimSuffix(hostStr, "/"))
	}
	return adaptedRoutesMap
}

// ReverseProxyHandler creates a reverse proxy for the specified target URL.
// It also supports optional path rewriting.
func ReverseProxyHandler(target string, rewritePath string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		if rewritePath != "" && strings.HasPrefix(req.URL.Path, rewritePath) {
			newPath := strings.TrimPrefix(req.URL.Path, rewritePath)
			if !strings.HasPrefix(newPath, "/") {
				newPath = "/" + newPath
			}
			req.URL.Path = newPath
		}
		req.Host = targetURL.Host
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying request to target: %s%s", target, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})
}

// Start starts the reverse proxy on the specified port
func (r *ReverseProxy) Start(port int) error {

	for pathPrefix, targetHost := range r.pathPrefixRoutesMap {
		http.Handle(string(pathPrefix), r.applyMiddlewares(ReverseProxyHandler(string(targetHost), string(pathPrefix))))
	}

	log.Printf("Starting reverse proxy server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
