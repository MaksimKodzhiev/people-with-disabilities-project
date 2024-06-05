package database

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Database struct {
    db                *sql.DB
    connectionTimeout time.Duration
    operationsTimeout time.Duration
    ParentContext     context.Context
}

func (database *Database) openDatabase() (err error) {
    database.connectionTimeout, err = time.ParseDuration(
        getEnv(
            "DATABASE_CONNECTION_TIMEOUT",
            defaultDatabaseConnectionTimeout,
        ),
    )

    if err != nil {
        return fmt.Errorf("openDatabase: could not parse database connection timeout: time.ParseDuration: %w", err)
    }

    database.operationsTimeout, err = time.ParseDuration(
        getEnv(
            "DATABASE_OPERATIONS_TIMEOUT",
            defaultDatabaseOperationsTimeout,
        ),
    )

    if err != nil {
        return fmt.Errorf("openDatabase: could not parse database operations timeout: time.ParseDuration: %w", err)
    }

    database.db, err = sql.Open(
        "mysql",
        fmt.Sprintf(
            "%s:%s@tcp(%s:%s)/project",
            getEnv("DATABASE_USERNAME", defaultDatabaseUsername),
            getEnv("DATABASE_PASSWORD", defaultDatabasePassword),
            getEnv("DATABASE_HOST", defaultDatabaseHost),
            getEnv("DATABASE_PORT", defaultDatabasePort),
        ),
    )

    if err != nil {
        return fmt.Errorf("openDatabase: could not open database: sql.Open: %w", err)
    }

    pingContext, pingCancelFunction := context.WithTimeout(
        database.ParentContext, database.connectionTimeout,
    )

    defer pingCancelFunction()

    err = database.db.PingContext(pingContext)

    if err != nil {
        return fmt.Errorf(
            "openDatabase: could not verify or establish database connection: sql.(*DB).PingContext: %w",
            err,
        )
    }

    return nil
}

func (database *Database) userExists(username string) (bool, error) {
    queryRowContext, queryRowCancelFunction := context.WithTimeout(
        database.ParentContext, database.operationsTimeout,
    )

    defer queryRowCancelFunction()

    err := database.db.QueryRowContext(
        queryRowContext, "SELECT username FROM users WHERE username=?", username,
    ).Scan(&username)

    return !errors.Is(err, sql.ErrNoRows), err
}

func (database *Database) createUser(username string, hashedPassword []byte) error {
    execContext, execCancelFunction := context.WithTimeout(
        database.ParentContext, database.operationsTimeout,
    )

    defer execCancelFunction()

    _, err := database.db.ExecContext(
        execContext, "INSERT INTO users (username, hashed_password) VALUES(?, ?)", username, hashedPassword,
    )

    if err != nil {
        return fmt.Errorf("createUser: sql.(*DB).ExecContext: %w", err)
    }

    return nil
}

func (database *Database) getHashedPassword(username string) ([]byte, error) {
    queryRowContext, queryRowCancelFunction := context.WithTimeout(
        database.ParentContext, database.operationsTimeout,
    )

    defer queryRowCancelFunction()

    var hashedPassword []byte

    err := database.db.QueryRowContext(
        queryRowContext, "SELECT hashed_password FROM users WHERE username=?", username,
    ).Scan(&hashedPassword)

    if err != nil {
        return nil, fmt.Errorf("getHashedPassword: sql.(*Row).Scan: %w", err)
    }

    return hashedPassword, nil
}
