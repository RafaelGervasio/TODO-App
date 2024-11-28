package main 

import (
    "TODO-App/router"
    "log"
    "TODO-App/middleware"
    "github.com/joho/godotenv"
    "os"
)


func main() {	
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}


    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        log.Fatal("DB_NAME is not set in the environment")
    }

	dbConn, err := middleware.InitDB(dbName)
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer dbConn.Close()


	router.StartRouter(dbConn)
}

