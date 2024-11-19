package middleware

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
    "fmt"
)

// InitDB initializes the database connection
func InitDB(dbName string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    fmt.Println("Database connected successfully")
    return db, nil
}
