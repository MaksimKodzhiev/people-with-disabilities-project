package middlewarepkg

import "net/http"

type WrappedResponseWriter struct {
    http.ResponseWriter
    httpStatusCode int
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
    return &WrappedResponseWriter{
        ResponseWriter: rw,
        httpStatusCode: http.StatusOK,
    }
}

func (w *WrappedResponseWriter) WriteHeader(httpStatusCode int) {
    w.httpStatusCode = httpStatusCode

    if w.httpStatusCode == http.StatusServiceUnavailable {
        w.Header().Set("Content-Type", "application/json")
    }

    w.ResponseWriter.WriteHeader(w.httpStatusCode)
}
