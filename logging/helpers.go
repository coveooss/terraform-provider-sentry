package logging

import (
	"encoding/json"
	"net/http"
	"strings"
)

// LogHttpResponse is a helper that logs at the specified logLevel the content of the HTTP
// response and the data that was extracted from it. It also includes some of the data
// contained in the request sent before getting the response
func LogHttpResponse(response *http.Response, responseData interface{}, logLevel LogLevel) {
	if auth, found := response.Request.Header["Authorization"]; found {
		// Scrub the auth header so no token leak during debug
		// Note: we are modifying the actual object, so the auth is permanently lost
		response.Request.Header["Authorization"] = []string{"Redacted to prevent leaks"}
		defer func() {
			response.Request.Header["Authorization"] = auth
		}()
	}

	prefix := strings.Repeat(" ", 6) // this is to respect the format string
	indent := strings.Repeat(" ", 2) // this is to pretty print
	jsonResponseDataBytes, jerr := json.MarshalIndent(responseData, prefix, indent)
	// Ignore err here since a map string => []string should be marshable #famouslastwords
	jsonRequestHeaders, _ := json.MarshalIndent(response.Request.Header, prefix, indent)
	jsonResponseHeaders, _ := json.MarshalIndent(response.Header, prefix, indent)

	// Use of the "." or else, first line is not formatted the right way
	logFormatString := `.
[HTTP REQUEST]
    Request URL: %s
    Request Headers:
      %v
[HTTP RESPONSE]
    Response Status: %s
    Response Headers: 
      %v
    Response Data:
      %+v
`
	logArgs := []interface{}{
		response.Request.URL.String(),
		string(jsonRequestHeaders),
		response.Status,
		string(jsonResponseHeaders),
	}
	if jerr == nil {
		logArgs = append(logArgs, string(jsonResponseDataBytes))
	} else {
		logArgs = append(logArgs, responseData)
	}

	logMethod := getLoggingFuncByLevel(logLevel)
	logMethod(logFormatString, logArgs...)
}
