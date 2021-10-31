# temaki


[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

![alt text](assets/temaki_logo.png?raw=true)

Minimal HTTP router based on the net/http standard library.

## Usage

Creating a _router_ object with **temaki** is very simple:

```golang
package main

import (
    "github.com/gyozatech/temaki"
    "log"
	"net/http"
)

func main() {
    router := temaki.NewRouter()
    
    router.UseMiddleware(exampleMiddleware)
    router.GET("api/v1/store/{storeId}/product/{productId}", getProductHandler)
    router.POST("api/v1/store/{storeId([^/]+)}/product/{productId([0-9]+)}", addProductHandler)
    router.PATCH("api/v1/store/{storeId([^/]+)}/product/{productId}", updateProductHandler)
    router.DELETE("api/v1/store/{storeId}/product/{productId([0-9]+)}", deleteProductHandler)

    log.Error(router.Start(8080))

}

// MIDDLEWARE
func exampleMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middleware before request phase!")
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
