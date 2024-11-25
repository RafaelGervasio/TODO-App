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
	"strings"
	"strconv"
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


func GetTasksHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.AuthenticateRequest(r, dbConn)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := getTasks(user.UserID, dbConn)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}


func GetTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.AuthenticateRequest(r, dbConn)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	taskID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := getTask(taskID, user.UserID, dbConn)
	if err != nil {
		http.Error(w, "Error fetching task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.AuthenticateRequest(r, dbConn)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var task models.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	task.UserID = user.UserID
	newTask, err := createTask(task, dbConn)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.AuthenticateRequest(r, dbConn)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	taskID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	task.TaskID = taskID
	task.UserID = user.UserID
	updatedTask, err := updateTask(task, dbConn)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := middleware.AuthenticateRequest(r, dbConn)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	taskID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := deleteTask(taskID, user.UserID, dbConn); err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getTasks(userID int, dbConn *sql.DB) ([]models.Task, error) {
	query := `SELECT task_id, user_id, name, description, due_date FROM tasks WHERE user_id = ?`
	rows, err := dbConn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.TaskID, &task.UserID, &task.Name, &task.Description, &task.DueDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func getTask(taskID, userID int, dbConn *sql.DB) (models.Task, error) {
	var task models.Task
	query := `SELECT task_id, user_id, name, description, due_date FROM tasks WHERE task_id = ? AND user_id = ?`
	err := dbConn.QueryRow(query, taskID, userID).Scan(&task.TaskID, &task.UserID, &task.Name, &task.Description, &task.DueDate)
	return task, err
}

func createTask(task models.Task, dbConn *sql.DB) (models.Task, error) {
	query := `INSERT INTO tasks (user_id, name, description, due_date) VALUES (?, ?, ?, ?)`
	result, err := dbConn.Exec(query, task.UserID, task.Name, task.Description, task.DueDate)
	if err != nil {
		return task, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return task, err
	}

	task.TaskID = int(taskID)
	return task, nil
}

func updateTask(task models.Task, dbConn *sql.DB) (models.Task, error) {
	query := `UPDATE tasks SET name = ?, description = ?, due_date = ? WHERE task_id = ? AND user_id = ?`
	_, err := dbConn.Exec(query, task.Name, task.Description, task.DueDate, task.TaskID, task.UserID)
	if err != nil {
		return task, err
	}

	return task, nil
}

func deleteTask(taskID, userID int, dbConn *sql.DB) error {
	query := `DELETE FROM tasks WHERE task_id = ? AND user_id = ?`
	_, err := dbConn.Exec(query, taskID, userID)
	return err
}



