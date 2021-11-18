package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var rgxHost = regexp.MustCompile(`\((.*?)\)`)

type RequestModifier func(req *http.Request)
type ResponseModifier func(*http.Response) error
type ErrorHandler func() func(http.ResponseWriter, *http.Request, error)
type Middleware func(handler http.Handler) http.Handler

func (rp *ReverseProxy) NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	if rp.ReqModifier != nil {
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			(*rp.ReqModifier)(req)
		}
	}
	if rp.RespModifier != nil {
		proxy.ModifyResponse = *rp.RespModifier
	}
	if rp.ErrHandler != nil {
		var errHandler func() func(http.ResponseWriter, *http.Request, error) = *rp.ErrHandler
		proxy.ErrorHandler = errHandler()
	}
	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

type ReverseProxy struct {
	Middlewares  *[]Middleware
	ReqModifier  *RequestModifier
	RespModifier *ResponseModifier
	ErrHandler   *ErrorHandler
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{&[]Middleware{}, nil, nil, nil}
}

func (rp *ReverseProxy) UseMiddleware(middleware ...Middleware) *ReverseProxy {
	*rp.Middlewares = append(*rp.Middlewares, middleware...)
	return rp
}

func (rp *ReverseProxy) UseRequestModifier(requestModifier RequestModifier) *ReverseProxy {
	rp.ReqModifier = &requestModifier
	return rp
}

func (rp *ReverseProxy) UseResponseModifier(responseModifier ResponseModifier) *ReverseProxy {
	rp.RespModifier = &responseModifier
	return rp
}

func (rp *ReverseProxy) UseErrorHandler(errorHandler ErrorHandler) *ReverseProxy {
	rp.ErrHandler = &errorHandler
	return rp
}

func (rp *ReverseProxy) Start(port int) error {
	rp.scanDomains()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (rp *ReverseProxy) scanDomains() map[string]*httputil.ReverseProxy {
	proxies := map[string]*httputil.ReverseProxy{}
	var envVars []string = os.Environ()
	for _, envVar := range envVars {
		envVarName := strings.Split(envVar, "=")[0]
		if strings.HasSuffix(envVarName, "_PROXY_URL") {
			rp.registerDomain(&proxies, envVarName)
		}
	}
	return proxies
}

func (rp *ReverseProxy) registerDomain(proxies *map[string]*httputil.ReverseProxy, envVarName string) {
	envVarValue := os.Getenv(envVarName)
	matches := rgxHost.FindAllString(envVarValue, -1)
	if len(matches) != 1 {
		panic(fmt.Errorf("invalid param passed for env var %s", envVarName))
	}
	host := matches[0]
	host = host[1 : len(host)-1]
	basePath := adaptPath(envVarValue)

	proxy, err := rp.NewProxy(host)
	if err != nil {
		panic(err)
	}
	var handler http.Handler = http.HandlerFunc(ProxyRequestHandler(proxy))
	for _, middleware := range *rp.Middlewares {
		handler = middleware(handler)
	}
	http.Handle(basePath, handler)
}

func adaptPath(envVarValue string) string {
	basePath := strings.Split(envVarValue, ")")[1]

	if basePath[0] != '/' {
		basePath = fmt.Sprintf("/%s", basePath)
	}
	if basePath[len(basePath)-1] != '/' {
		basePath = fmt.Sprintf("%s/", basePath)
	}
	return basePath
}
