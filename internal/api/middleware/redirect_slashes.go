package middlewarepkg

import (
    "net/http"
    "strings"
)

func RedirectSlashes(next http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        path := r.URL.Path

        if len(path) > 1 && strings.HasSuffix(path, "/") {
            path = strings.TrimSuffix(path, "/")

            if r.URL.RawQuery != "" {
                path += "?" + r.URL.RawQuery
            }

            var httpStatusCode int

            if r.Method == http.MethodGet {
                httpStatusCode = http.StatusMovedPermanently
            } else {
                httpStatusCode = http.StatusPermanentRedirect
            }

            http.Redirect(rw, r, path, httpStatusCode)

            return
        }

        next(rw, r)
    }
}
