package main

import (
	"auth_api/controller"
	"auth_api/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error
	client := db.ConnectDB()
	userController := controller.NewUserController(client)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			userController.Login(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			userController.Register(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is starting at http://localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
	fmt.Println("Server started successfully!")
}
