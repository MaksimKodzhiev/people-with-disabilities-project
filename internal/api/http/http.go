package http

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/iancoleman/orderedmap"
)

var httpStatusCodes = map[int]string{
    http.StatusContinue:           "Continue",            // RFC 9110, 15.2.1
    http.StatusSwitchingProtocols: "Switching Protocols", // RFC 9110, 15.2.2
    http.StatusProcessing:         "Processing",          // RFC 2518, 10.1
    http.StatusEarlyHints:         "Early Hints",         // RFC 8297

    http.StatusOK:                   "OK",                            // RFC 9110, 15.3.1
    http.StatusCreated:              "Created",                       // RFC 9110, 15.3.2
    http.StatusAccepted:             "Accepted",                      // RFC 9110, 15.3.3
    http.StatusNonAuthoritativeInfo: "Non-Authoritative Information", // RFC 9110, 15.3.4
    http.StatusNoContent:            "No Content",                    // RFC 9110, 15.3.5
    http.StatusResetContent:         "Reset Content",                 // RFC 9110, 15.3.6
    http.StatusPartialContent:       "Partial Content",               // RFC 9110, 15.3.7
    http.StatusMultiStatus:          "Multi-Status",                  // RFC 4918, 11.1
    http.StatusAlreadyReported:      "Already Reported",              // RFC 5842, 7.1
    http.StatusIMUsed:               "IM Used",                       // RFC 3229, 10.4.1

    http.StatusMultipleChoices:   "Multiple Choices",   // RFC 9110, 15.4.1
    http.StatusMovedPermanently:  "Moved Permanently",  // RFC 9110, 15.4.2
    http.StatusFound:             "Found",              // RFC 9110, 15.4.3
    http.StatusSeeOther:          "See Other",          // RFC 9110, 15.4.4
    http.StatusNotModified:       "Not Modified",       // RFC 9110, 15.4.5
    http.StatusUseProxy:          "Use Proxy",          // RFC 9110, 15.4.6
    http.StatusTemporaryRedirect: "Temporary Redirect", // RFC 9110, 15.4.8
    http.StatusPermanentRedirect: "Permanent Redirect", // RFC 9110, 15.4.9

    http.StatusBadRequest:                   "Bad Request",                     // RFC 9110, 15.5.1
    http.StatusUnauthorized:                 "Unauthorized",                    // RFC 9110, 15.5.2
    http.StatusPaymentRequired:              "Payment Required",                // RFC 9110, 15.5.3
    http.StatusForbidden:                    "Forbidden",                       // RFC 9110, 15.5.4
    http.StatusNotFound:                     "Not Found",                       // RFC 9110, 15.5.5
    http.StatusMethodNotAllowed:             "Method Not Allowed",              // RFC 9110, 15.5.6
    http.StatusNotAcceptable:                "Not Acceptable",                  // RFC 9110, 15.5.7
    http.StatusProxyAuthRequired:            "Proxy Authentication Required",   // RFC 9110, 15.5.8
    http.StatusRequestTimeout:               "Request Timeout",                 // RFC 9110, 15.5.9
    http.StatusConflict:                     "Conflict",                        // RFC 9110, 15.5.10
    http.StatusGone:                         "Gone",                            // RFC 9110, 15.5.11
    http.StatusLengthRequired:               "Length Required",                 // RFC 9110, 15.5.12
    http.StatusPreconditionFailed:           "Precondition Failed",             // RFC 9110, 15.5.13
    http.StatusRequestEntityTooLarge:        "Request Entity Too Large",        // RFC 9110, 15.5.14
    http.StatusRequestURITooLong:            "Request URI Too Long",            // RFC 9110, 15.5.15
    http.StatusUnsupportedMediaType:         "Unsupported Media Type",          // RFC 9110, 15.5.16
    http.StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable", // RFC 9110, 15.5.17
    http.StatusExpectationFailed:            "Expectation Failed",              // RFC 9110, 15.5.18
    http.StatusTeapot:                       "I'm a teapot",                    // RFC 9110, 15.5.19 (Unused)
    http.StatusMisdirectedRequest:           "Misdirected Request",             // RFC 9110, 15.5.20
    http.StatusUnprocessableEntity:          "Unprocessable Entity",            // RFC 9110, 15.5.21
    http.StatusLocked:                       "Locked",                          // RFC 4918, 11.3
    http.StatusFailedDependency:             "Failed Dependency",               // RFC 4918, 11.4
    http.StatusTooEarly:                     "Too Early",                       // RFC 8470, 5.2
    http.StatusUpgradeRequired:              "Upgrade Required",                // RFC 9110, 15.5.22
    http.StatusPreconditionRequired:         "Precondition Required",           // RFC 6585, 3
    http.StatusTooManyRequests:              "Too Many Requests",               // RFC 6585, 4
    http.StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large", // RFC 6585, 5
    http.StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",   // RFC 7725, 3

    http.StatusInternalServerError:           "Internal Server Error",           // RFC 9110, 15.6.1
    http.StatusNotImplemented:                "Not Implemented",                 // RFC 9110, 15.6.2
    http.StatusBadGateway:                    "Bad Gateway",                     // RFC 9110, 15.6.3
    http.StatusServiceUnavailable:            "Service Unavailable",             // RFC 9110, 15.6.4
    http.StatusGatewayTimeout:                "Gateway Timeout",                 // RFC 9110, 15.6.5
    http.StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",      // RFC 9110, 15.6.6
    http.StatusVariantAlsoNegotiates:         "Variant Also Negotiates",         // RFC 2295, 8.1
    http.StatusInsufficientStorage:           "Insufficient Storage",            // RFC 4918, 11.5
    http.StatusLoopDetected:                  "Loop Detected",                   // RFC 5842, 7.2
    http.StatusNotExtended:                   "Not Extended",                    // RFC 2774, 7
    http.StatusNetworkAuthenticationRequired: "Network Authentication Required", // RFC 6585, 6
}

func writeJSONResponse(rw http.ResponseWriter, httpStatusCode int, response any) error {
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(httpStatusCode)

    err := json.NewEncoder(rw).Encode(response)

    if err != nil {
        return fmt.Errorf("writeJSONResponse: json.(*Encoder).Encode: %w", err)
    }

    return nil
}

func WriteJSONResponse(rw http.ResponseWriter, httpStatusCode int, result any) error {
    orderedMap := orderedmap.New()

    orderedMap.Set("ok", true)
    orderedMap.Set("result", result)

    err := writeJSONResponse(rw, httpStatusCode, orderedMap)

    if err != nil {
        return fmt.Errorf("WriteJSONResponse: %w", err)
    }

    return nil
}

func WriteErrorJSONResponse(rw http.ResponseWriter, r *http.Request, httpStatusCode int, message string) error {
    httpStatusCodeDescription, ok := httpStatusCodes[httpStatusCode]

    if !ok {
        return fmt.Errorf("WriteErrorJSONResponse: http status code %d is not supported", httpStatusCode)
    }

    orderedMap := orderedmap.New()

    orderedMap.Set("ok", false)
    orderedMap.Set("error_code", httpStatusCode)
    orderedMap.Set("description", httpStatusCodeDescription)
    orderedMap.Set("timestamp", time.Now().Format(time.RFC3339))
    orderedMap.Set("path", r.URL.Path)
    orderedMap.Set("message", message)

    err := writeJSONResponse(rw, httpStatusCode, orderedMap)

    if err != nil {
        return fmt.Errorf("WriteErrorJSONResponse: %w", err)
    }

    return nil
}

func WriteGenericErrorJSONResponse(rw http.ResponseWriter, r *http.Request, httpStatusCode int) error {
    httpStatusCodeDescription, ok := httpStatusCodes[httpStatusCode]

    if !ok {
        return fmt.Errorf("WriteGenericErrorJSONResponse: http status code %d is not supported", httpStatusCode)
    }

    orderedMap := orderedmap.New()

    orderedMap.Set("ok", false)
    orderedMap.Set("error_code", httpStatusCode)
    orderedMap.Set("description", httpStatusCodeDescription)
    orderedMap.Set("timestamp", time.Now().Format(time.RFC3339))
    orderedMap.Set("path", r.URL.Path)

    err := writeJSONResponse(rw, httpStatusCode, orderedMap)

    if err != nil {
        return fmt.Errorf("WriteGenericErrorJSONResponse: %w", err)
    }

    return nil
}
