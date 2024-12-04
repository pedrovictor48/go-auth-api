package main

import (
	"auth_api/controller"
	"auth_api/db"
	"auth_api/repository"
	"auth_api/usecase"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error
	client := db.ConnectDB()
	userRepository := repository.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)

	postRepository := repository.NewPostRepository(client)
	postUsecase := usecase.NewPostUsecase(postRepository)
	postController := controller.NewPostController(postUsecase)

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			postController.CreatePost(w, r)
		} else if http.MethodGet == r.Method {
			postController.GetPostsById(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

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
