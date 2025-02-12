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
func (rp *ReverseProxy) WithMiddlewares(middlewares ...Middleware) *ReverseProxy {
	if rp.middlewares == nil {
		rp.middlewares = []Middleware{}
	}
	rp.middlewares = append(rp.middlewares, middlewares...)
	return rp
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

func (rp *ReverseProxy) applyMiddlewares(handler http.Handler) http.Handler {
	for _, m := range rp.middlewares {
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
		if !strings.HasPrefix(hostStr, "http") && !strings.HasPrefix(hostStr, "ws") {
			// setting http as default protocol if no https/http or wss/ws protocol is specified
			hostStr = "http://" + hostStr
		}
		adaptedRoutesMap[PathPrefix(prefixStr)] = TargetHost(strings.TrimSuffix(hostStr, "/"))
	}
	return adaptedRoutesMap
}

// Start starts the reverse proxy on the specified port
func (rp *ReverseProxy) Start(port int) error {
	// we manage every prefix from a single / root path to avoid inconvenient HTTP statuses 302
	http.Handle("/", rp.applyMiddlewares(toHTTPHandler(rp.handleFunc)))

	log.Printf("Starting reverse proxy server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (rp *ReverseProxy) handleFunc(w http.ResponseWriter, req *http.Request) {
	for pathPrefix, targetHost := range rp.pathPrefixRoutesMap {
		if strings.HasPrefix(req.URL.Path, strings.TrimSuffix(string(pathPrefix), "/")) {
			targetURL, err := url.Parse(string(targetHost))
			if err != nil {
				log.Printf("Error parsing target URL: %v", err)
				http.Error(w, fmt.Sprintf("Error parsing target URL: %s", err), http.StatusInternalServerError)
				return
			}

			proxy := httputil.NewSingleHostReverseProxy(targetURL)

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				req.Header.Set("X-Forwarded-Host", req.Host)
				req.Header.Set("X-Real-IP", req.RemoteAddr)

				rewritePath := string(pathPrefix)
				if rewritePath != "" && strings.HasPrefix(req.URL.Path, rewritePath) {
					newPath := strings.TrimPrefix(req.URL.Path, rewritePath)
					if !strings.HasPrefix(newPath, "/") {
						newPath = "/" + newPath
					}
					req.URL.Path = newPath
				}
				req.Host = targetURL.Host
			}

			// WebSocket request
			if isWebSocketRequest(req) {
				log.Println("Handling WebSocket Upgrade...")
				handleWebSocket(w, req, string(targetHost))
				return
			}

			// HTTP/HTTPS request
			log.Printf("Proxying request to target: %s%s", targetHost, req.URL.Path)
			proxy.ServeHTTP(w, req)
			return
		}
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

// convert a func(http.ResponseWriter, *http.Request) into an http.Handler to adapt http.HandleFunc to http.Handle
func toHTTPHandler(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(fn)
}
