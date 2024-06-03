package constants

import "time"

const (
    ReadTimeout    = 5 * time.Second
    WriteTimeout   = 5 * time.Second
    RequestTimeout = 5 * time.Second
    IdleTimeout    = 30 * time.Second
)

const RequestTimeoutMessage = `{"ok":false,"error_code":503,"description":"Service Unavailable"}`
