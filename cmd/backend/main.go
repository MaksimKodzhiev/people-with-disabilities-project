package main

import (
    "log"

    project "project/internal"
)

func main() {
    err, shutdownErr := project.RunHTTPServer()

    if err == nil && shutdownErr == nil {
        log.Println("Server terminated successfully!")
    } else {
        if shutdownErr != nil {
            log.Println(shutdownErr)
        }

        if err != nil {
            log.Println(err)
        }
    }
}
