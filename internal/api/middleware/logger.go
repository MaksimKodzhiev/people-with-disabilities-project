package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        wrappedResponseWriter, ok := rw.(*WrappedResponseWriter)

        start := time.Now()

        if ok {
            next(wrappedResponseWriter, r)

            log.Println(r.RemoteAddr, wrappedResponseWriter.httpStatusCode, r.Method, r.RequestURI, time.Since(start))
        } else {
            next(rw, r)

            log.Println(r.RemoteAddr, r.Method, r.RequestURI, time.Since(start))
        }
    }
}
