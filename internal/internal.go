package internal

import (
    "context"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    . "project/internal/api/constants"
    "project/internal/api/middleware"
)

func RunHTTPServer() (err, shutdownErr error) {
    middlewareChain := middlewarepkg.Chain(
        middlewarepkg.WrapResponseWriter,
        middlewarepkg.Logger,
        middlewarepkg.Timeout,
        middlewarepkg.Recoverer,
        middlewarepkg.RedirectSlashes,
    )

    serverMultiplexer := http.NewServeMux()

    server := http.Server{
        Addr:         ":http",
        ReadTimeout:  ReadTimeout,
        WriteTimeout: WriteTimeout + RequestTimeout,
        IdleTimeout:  IdleTimeout,
        Handler:      middlewareChain(serverMultiplexer.ServeHTTP),
    }

    go func() {
        log.Println("Server is listening on", server.Addr)

        err = server.ListenAndServe()

        if !errors.Is(err, http.ErrServerClosed) {
            err = fmt.Errorf("RunHTTPServer: http.(*Server).ListenAndServe: %w", err)
        } else {
            err = nil
        }
    }()

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
    <-signalChannel

    ctx, cancel := context.WithTimeout(context.Background(), WriteTimeout+RequestTimeout)

    defer cancel()

    shutdownErr = server.Shutdown(ctx)

    if shutdownErr != nil {
        shutdownErr = fmt.Errorf("RunHTTPServer: http.(*Server).Shutdown: %w", shutdownErr)
    }

    return err, shutdownErr
}
