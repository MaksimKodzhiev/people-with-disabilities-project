package middlewarepkg

import "net/http"

func WrapResponseWriter(next http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        next(NewWrappedResponseWriter(rw), r)
    }
}
