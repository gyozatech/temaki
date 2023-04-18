# temaki


[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

![alt text](assets/logo.png?raw=true)

Minimal HTTP router based on the net/http standard library.

## Usage

Creating a _router_ object with **temaki** is very simple:

```golang
package main

import (
    "github.com/gyozatech/temaki"
    "github.com/gyozatech/temaki/middlewares"
    "log"
    "net/http"
    "fmt"
    "strconv"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func main() {
    router := temaki.NewRouter()

    // passing a custom logger to the middlewares (the default is gyozatch/noodlog)
    middlewares.SetLogger(logger)
    
    // provided temaki middlewares
    router.UseMiddleware(middlewares.RecoverPanicMiddleware)
    router.UseMiddleware(middlewares.RequestLoggerMiddleware)
    router.UseMiddleware(middlewares.CORSMiddleware)
    // custom middleware
    router.UseMiddleware(authMiddleware)

    // routes
    router.GET("/api/v1/stores/{storeId}/products/{productId}", getProductHandler)
    router.POST("/api/v1/stores/{storeId([^/]+)}/products/{productId([0-9]+)}", addProductHandler)
    router.PATCH("/api/v1/stores/{storeId([^/]+)}/products/{productId}", updateProductHandler)
    router.DELETE("/api/v1/stores/{storeId}/products/{productId([0-9]+)}", deleteProductHandler)

    log.Fatal(router.Start(8080))
}

// MIDDLEWARE
func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middleware before request phase!")
		token, err := temaki.GetBearerToken(r)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if token != "IcgFd3FiHwKDM2H" {
            http.Error(w, "invalid bearer token provided", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
		fmt.Println("Executing middleware after request phase!")
	})
}

// HANDLERS ********

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	storeId := temaki.GetPathParam(r, "storeId")
	productId, _ := strconv.Atoi(temaki.GetPathParam(r, "productId"))
	fmt.Fprintf(w, "getProductHandler %s %d\n", storeId, productId)
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	storeId := temaki.GetPathParam(r, "storeId")
	productId, _ := strconv.Atoi(temaki.GetPathParam(r, "productId"))
	fmt.Fprintf(w, "addProductHandler %s %d\n", storeId, productId)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	storeId := temaki.GetPathParam(r, "storeId")
	productId, _ := strconv.Atoi(temaki.GetPathParam(r, "productId"))
	fmt.Fprintf(w, "updateProductHandler %s %d\n", storeId, productId)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	storeId := temaki.GetPathParam(r, "storeId")
	productId, _ := strconv.Atoi(temaki.GetPathParam(r, "productId"))
	fmt.Fprintf(w, "deleteProductHandler %s %d\n", storeId, productId)
}

```

As you have noticed from the _paths_ you can also decide to specify a **regex** pattern to the path parameters.

## Contributing

Any contribution to this project is welcome! Just fork the project, and open a Pull Request.

If you find a problem or a bug, or have some improvement suggestion, please open an issue or a discussion on Github.

## Special thanks

Thanks to [benhoyt](https://github.com/benhoyt/go-routing) who gaves me the idea for this wrapper.
