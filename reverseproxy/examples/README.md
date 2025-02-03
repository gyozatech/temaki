# Example of usage of the reverseproxy

You need to download the current directory and run each service in a different tab.
You must start with `main_s1.go` and `main_s2.go` representing the two backend services to which the reverseproxy is going to proxy the requests which is receiving.

In the example we're using env vars to register the routes toward the backend services `s1` and `s2` based on the prefix in the HTTP request.

## Usage  
  
After having run:
- `go run main_s1.go` on one tab
- `go run main_s2.go` on another tab

You can run the reverse proxy via:
- `go run main.go` onto another tab

Then, you can test the routes with the following endpoints:
```
curl -H "Authorization: Bearer abcd" http://localhost:8080/s1/api/v1/status
curl -H "Authorization: Bearer abcd" http://localhost:8080/s1/api/v1/hello
curl -H "Authorization: Bearer abcd" http://localhost:8080/s2/api/v1/status
curl -H "Authorization: Bearer abcd" http://localhost:8080/s2/api/v1/hello
```
