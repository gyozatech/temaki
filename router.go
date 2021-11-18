package temaki

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	Routes      *[]Route
	Middlewares *[]Middleware
}

func NewRouter() *Router {
	return &Router{&[]Route{}, &[]Middleware{}}
}

func (router *Router) UseMiddleware(middleware ...Middleware) *Router {
	*router.Middlewares = append(*router.Middlewares, middleware...)
	return router
}

func (router *Router) GET(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("GET", pattern, handlerFunc))
	return router
}

func (router *Router) POST(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("OPTIONS", pattern, handlerFunc))
	*router.Routes = append(*router.Routes, NewRoute("POST", pattern, handlerFunc))
	return router
}

func (router *Router) PUT(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("OPTIONS", pattern, handlerFunc))
	*router.Routes = append(*router.Routes, NewRoute("PUT", pattern, handlerFunc))
	return router
}

func (router *Router) PATCH(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("OPTIONS", pattern, handlerFunc))
	*router.Routes = append(*router.Routes, NewRoute("PATCH", pattern, handlerFunc))
	return router
}

func (router *Router) DELETE(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("OPTIONS", pattern, handlerFunc))
	*router.Routes = append(*router.Routes, NewRoute("DELETE", pattern, handlerFunc))
	return router
}

func (router *Router) OPTIONS(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("OPTIONS", pattern, handlerFunc))
	return router
}

func (router *Router) HEAD(pattern string, handlerFunc http.HandlerFunc) *Router {
	*router.Routes = append(*router.Routes, NewRoute("HEAD", pattern, handlerFunc))
	return router
}

func (router *Router) DispatcherHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allow []string
		for _, route := range *router.Routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				if r.Method != route.method {
					allow = append(allow, route.method)
					continue
				}

				r = enrichRequestContext(r, ctxKey{}, matches[1:])
				r = enrichRequestContext(r, "pathParamsMap", route.pathParams)

				route.handler(w, r)
				return
			}
		}
		if len(allow) > 0 {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, r)
	}
}

func (router *Router) Serve() http.Handler {
	var finalHandler http.Handler = http.HandlerFunc(router.DispatcherHandler())

	for _, middleware := range *router.Middlewares {
		finalHandler = middleware(finalHandler)
	}
	return finalHandler
}

func (router *Router) Start(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router.Serve())
}

func enrichRequestContext(r *http.Request, key, val interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}
