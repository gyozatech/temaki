# Reverse Proxy

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

This package contains a simple HTTP/S Reverse Proxy meant to be imported in your application and configured programmatically and through environment variables.

## How it works

You must list your services though enviroment variable by using the following convention:

name: `<SERVICENAME>_PROXY_URL`
value: `(HOST)/common/path/`
  
to get all the calls coming to the reverse proxy host at the path `/common/path/*` to be proxied to their respective hosts.
  
For example if you have your reverse proxy server listening to `http://localhost:8080` and your backend service listening to `http://localhost:8081/service-1/api/v1/*`

You can get the calls coming to your reverse proxy host at the paths `http://localhost:8080/service-1/api/v1/*` to be all proxied to the backend service at the paths `http://localhost:8081/service-1/api/v1/*`.

This mechanism can be created by setting the following environment variable:

```bash
export SERVICE_1_PROXY_URL='(http://localhost:8081)/service-1/api/v1/'
```
And by programming your reverse proxy this way:

```go
package main

import (
	"log"
	"github.com/gyozatech/temaki/reverseproxy"
)

func main() {
	log.Fatal(reverseproxy.NewReverseProxy().Start(8080))
}
```
