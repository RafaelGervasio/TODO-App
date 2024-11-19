package models

type User struct {
    UserID   int    `json:"user_id" db:"user_id"`  // Maps to `user_id` in the database
    Name     string `json:"name" db:"name"`       // Maps to `name` in the database
    Email    string `json:"email" db:"email"`     // Maps to `email` in the database
    Password string `json:"password" db:"password"` // Maps to `password` in the database
}

type Task struct {
    TaskID      int    `json:"task_id" db:"task_id"`        // Maps to `task_id` in the database
    UserID      int    `json:"user_id" db:"user_id"`        // Maps to `user_id` in the database (foreign key)
    Name        string `json:"name" db:"name"`             // Maps to `name` in the database
    Description string `json:"description,omitempty" db:"description"` // Maps to `description` in the database, omits empty JSON
    DueDate     string `json:"due_date,omitempty" db:"due_date"`       // Maps to `due_date` in the database, omits empty JSON
}
