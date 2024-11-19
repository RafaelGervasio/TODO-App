package handlers

import (
	"net/http"
 	"TODO-App/models"
	"TODO-App/middleware"
	"encoding/json"
	"fmt"
)

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

	if err := registerUser(user, dbConn *sql.DB); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}


func registerUser(user models.User, dbConn *sql.DB) error {
    query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
    _, err = dbConn.Exec(query, user.Name, user.Email, user.Password)
    return err
}



func LoginHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	// only allow post method
	// get the body - decode the json into a go struct
	// attempt to login - return error if not possible
	// return success


	if r.Method	!= http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
	}

	
}
