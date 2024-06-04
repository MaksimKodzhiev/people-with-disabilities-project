package database

import "os"

const (
    defaultDatabaseUsername = "admin"
    defaultDatabasePassword = "admin"
    defaultDatabaseHost     = "localhost"
    defaultDatabasePort     = "3306"
)

const (
    defaultDatabaseConnectionTimeout = "10s"
    defaultDatabaseOperationsTimeout = "5s"
)

func getEnv(key, fallback string) string {
    value, ok := os.LookupEnv(key)

    if ok {
        return value
    }

    return fallback
}
