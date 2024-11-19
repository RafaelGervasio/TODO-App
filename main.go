package main 

import (
    "TODO-App/router"
    "log"
    "TODO-App/middleware"
    // "database/sql"
    // _ "github.com/mattn/go-sqlite3" // Import SQLite driver
)


func main() {	
	dbConn, err := middleware.InitDB("todo.db")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer dbConn.Close()

	router.StartRouter(dbConn)
}

