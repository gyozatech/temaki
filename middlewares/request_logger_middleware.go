package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"time"
)

// RequestLoggerMiddleware is the middleware layer to log all the HTTP requests
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rww := NewResponseWriterWrapper(&w, r)

		defer func() {

			defer func() {
				if rec := recover(); rec != nil {
					logError("Panic recovered in the RequestLoggerMiddleware: ", string(debug.Stack()))
				}
			}()

			logInfo(ReqRespLogStruct{
				Request:  HTTPRequest(r),
				ExecTime: time.Since(start),
				Response: HTTPResponse(*rww.statusCode, rww.Header(), rww.r.RequestURI, rww.body.String()),
			})

		}()

		next.ServeHTTP(rww, r)

	})
}

// ResponseWriterWrapper struct is used to log the response
type ResponseWriterWrapper struct {
	r          *http.Request
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

// NewResponseWriterWrapper static function creates a wrapper for the http.ResponseWriter
func NewResponseWriterWrapper(w *http.ResponseWriter, r *http.Request) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		r:          r,
		w:          w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}

/////////////////////////////////////////////////////////////////////////////////////////////

// ReqRespLogStruct struct represents the schema for the HTTP Request/ExecutionTime/Response log
type ReqRespLogStruct struct {
	Request  RequestStruct
	ExecTime time.Duration
	Response ResponseStruct
}

// RequestStruct struct represents the schema for the HTTP Request log
type RequestStruct struct {
	Verb     string            `json:"verb,omitempty"`
	Path     string            `json:"path,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     interface{}       `json:"body,omitempty"`
}

// ResponseStruct struct represents the schema for the HTTP Response log
type ResponseStruct struct {
	Status  int               `json:"status,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
}

// HTTPRequest static function flatten pointers of HTTP Request and obscurate passwords to prepare for logging
func HTTPRequest(r *http.Request) RequestStruct {
	if r == nil {
		return RequestStruct{}
	}

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		logInfo(err)
	}

	var req RequestStruct

	lines := strings.Split(string(requestDump), "\r\n")
	dir := strings.Split(lines[0], " ")
	req.Verb = dir[0]
	req.Path = dir[1]
	req.Protocol = dir[2]

	var headers = map[string]string{}
	for _, h := range lines[1 : len(lines)-3] {
		head := strings.Split(h, ":")
		headers[head[0]] = strings.ReplaceAll(head[1], " ", "")
	}
	req.Headers = headers
	req.Body = getObjFromJSON([]byte(lines[len(lines)-1]))
	return req
}

// HTTPResponse static function flatten pointers of HTTP Response and obscurate passwords to prepare for logging
func HTTPResponse(statusCode int, responseHeaders http.Header, requestURI string, responseBody string) ResponseStruct {
	resp := ResponseStruct{Status: statusCode}
	headers := map[string]string{}
	for k, v := range responseHeaders {
		header := ""
		for _, hv := range v {
			header = header + hv + ";"
		}
		headers[k] = header[0 : len(header)-1]
	}
	resp.Headers = headers

	var bodyObj interface{}
	_ = json.Unmarshal([]byte(responseBody), &bodyObj)
	resp.Body = bodyObj
	return resp
}

func getObjFromJSON(jsonMsg []byte) interface{} {
	var obj interface{}
	_ = json.Unmarshal(jsonMsg, &obj)
	return obj
}
