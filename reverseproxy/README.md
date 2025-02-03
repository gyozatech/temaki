# Reverse Proxy

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

This package contains a simple HTTP/S Reverse Proxy meant to be imported in your application and configured programmatically and through environment variables.

## How it works

You must list your services through enviroment variable by using the following convention:

name: `PROXY_RULE_<SERVICENAME>`
value: `/<pathprefix>/>https://<host:port>`
  
The incoming requests having the specified path prefix will be completely redirected to the given host, removing the path prefix and appending the rest of the URL as long as any headers, cookies, query params.

Example:
```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gyozatech/temaki/reverseproxy"
	"github.com/gyozatech/temaki/middlewares"
)

func main() {

	os.Setenv("PROXY_RULE_S1", "/weather/>https://api.weather.com")
	os.Setenv("PROXY_RULE_S2", "/geo/>https://api.geo.com")

	routes := reverseproxy.CollectPathPrefixRoutesFromEnvVar()
	// or alternatively, you can initialize directly the PathPrefixRoutesMap: 
	/* 
           routes := reverseproxy.PathPrefixRoutesMap{
		"weather": "https://api.weather.com",
		"geo": "https://api.geo.com",
	}
        */

	log.Fatalf("Server error: %s", reverseproxy.New(routes).
		WithMiddlewares(middlewares.RequestLoggerMiddleware).
		Start(8080))
}
```
