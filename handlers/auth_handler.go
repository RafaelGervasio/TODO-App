package handlers

import (
	"net/http"
 	"TODO-App/models"
	"TODO-App/middleware"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt"
	"time"
	"database/sql"
)


var secretKey = []byte("secureSecretText")


func RegisterHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	// only allow post method
	// get the body - decode the json into a go struct
	// attempt to register in the database - return error if not possible
	// return success

	if r.Method	!= http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
        return
	}

	if err := registerUser(user, dbConn); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}


func registerUser(user models.User, dbConn *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err)
	}

    query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
    _, err = dbConn.Exec(query, user.Name, user.Email, string(hashedPassword))
    return err
}



func LoginHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	// only allow post method
	// get the body - decode the json into a go struct
	// attempt to login - return error if not possible
	// create JWT
	// return success
	
	if r.Method	!= http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
	}

	var credentials models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := authenticateUser(credentials.Email, credentials.Password, dbConn)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := middleware.CustomClaims{
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Set expiration time
			Issuer:    "TODO-App",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"token": signedToken,
	}

	json.NewEncoder(w).Encode(response)
}


func authenticateUser (email, password string, dbConn *sql.DB) (models.User, error) {
	var user models.User

	query := `SELECT user_id, name, email, password FROM users WHERE email = ?`
	err := dbConn.QueryRow(query, email).Scan(&user.UserID, &user.Name, &user.Email, &user.Password)
	
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	}
	if err != nil {
		return user, fmt.Errorf("error querying database: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("invalid credentials")
	}

	return user, nil

}


