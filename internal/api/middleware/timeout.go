package middleware

import (
    "net/http"

    . "project/internal/api/constants"
)

func Timeout(next http.HandlerFunc) http.HandlerFunc {
    return http.TimeoutHandler(next, WriteTimeout, RequestTimeoutMessage).ServeHTTP
}
