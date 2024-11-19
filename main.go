package main 

import (
    "TODO-App/router"
    // "database/sql"
    // _ "github.com/mattn/go-sqlite3" // Import SQLite driver
)


func main() {
	router.StartRouter()
}



// func initDB(dbName string) (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", dbName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
//         return nil, err
//     }

//     return db, nil
// }