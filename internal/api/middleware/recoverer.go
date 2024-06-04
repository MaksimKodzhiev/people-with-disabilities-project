package middleware

import (
    "log"
    "net/http"
    "runtime/debug"
)

func Recoverer(next http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()

            if err != nil {
                if err == http.ErrAbortHandler {
                    panic(err)
                } else {
                    log.Println(err, string(debug.Stack()))
                }
            }
        }()

        next(rw, r)
    }
}
