package middleware

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
    "fmt"
    "github.com/golang-jwt/jwt"
    "TODO-App/models"
    "strings"
    "net/http"
    "time"
    "github.com/joho/godotenv"
    "os"
    "log"
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


func GetJWTSecret() []byte {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    secretKey := os.Getenv("JWT_SECRET")
    if secretKey == "" {
        log.Fatal("JWT_SECRET is not set in the .env file")
    }
    
    return []byte(secretKey)
}



type CustomClaims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}


func AuthenticateRequest(r *http.Request, dbConn *sql.DB) (models.User, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return models.User{}, fmt.Errorf("authorization header required")
    }

    bearerToken := strings.Split(authHeader, "Bearer ")
    if len(bearerToken) != 2 {
        return models.User{}, fmt.Errorf("invalid token format")
    }

    tokenString := bearerToken[1]

    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return GetJWTSecret(), nil
    })

    if err != nil {
        return models.User{}, fmt.Errorf("invalid token: %v", err)
    }

    claims, ok := token.Claims.(*CustomClaims)
    if !ok {
        return models.User{}, fmt.Errorf("invalid token claims")
    }

    if claims.ExpiresAt < time.Now().UTC().Unix() {
        return models.User{}, fmt.Errorf("token has expired")
    }

    var user models.User
    query := `SELECT user_id, name, email FROM users WHERE name = ?`
    err = dbConn.QueryRow(query, claims.Username).Scan(&user.UserID, &user.Name, &user.Email)
    if err != nil {
        return models.User{}, fmt.Errorf("user not found")
    }

    return user, nil
}
