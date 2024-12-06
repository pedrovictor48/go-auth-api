package main

import (
	"auth_api/controller"
	"auth_api/db"
	"auth_api/repository"
	"auth_api/usecase"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	body               map[string]string
	expectedStatusCode int
}

func TestLogin(t *testing.T) {
	client := db.ConnectDB()
	userRepository := repository.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)

	tests := []TestCase{
		{
			body: map[string]string{
				"email":    "teste@gmail.com",
				"password": "123",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			body: map[string]string{
				"email":    "teste@gmail.com",
				"password": "wrongpassword",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			body: map[string]string{
				"email":    "test@gmail.com",
				"password": "123",
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		body := test.body
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyJSON))

		userController.Login(rr, req)

		if test.expectedStatusCode != rr.Code {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	}

}

func TestPost(t *testing.T) {
	client := db.ConnectDB()
	userRepository := repository.NewUserRepository(client)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controller.NewUserController(userUsecase)

	postRepository := repository.NewPostRepository(client)
	postUsecase := usecase.NewPostUsecase(postRepository)
	postController := controller.NewPostController(postUsecase)

	// login
	body := map[string]string{
		"email":    "teste@gmail.com",
		"password": "123",
	}

	bodyJSON, _ := json.Marshal(body)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyJSON))

	userController.Login(rr, req)

	tests := []TestCase{
		{
			body: map[string]string{
				"title":   "Title 1",
				"content": "Content 1",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			body: map[string]string{
				"title":   "Title 2",
				"content": "Content 2",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			body: map[string]string{
				"title":   "Title 3",
				"content": "Content 3",
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		body := test.body
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(bodyJSON))

		postController.CreatePost(rr, req)

		if test.expectedStatusCode != rr.Code {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	}
}
