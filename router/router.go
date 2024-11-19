package router

import (
	"fmt"
	"net/http"
	"TODO-App/handlers"
)

func StartRouter () {
	mux := http.NewServeMux()


	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "200 OK")
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/register", handlers.RegisterHandler)


	fmt.Println("Server running on http://localhost:8080")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting sercer: %v\n", err)
	}
}