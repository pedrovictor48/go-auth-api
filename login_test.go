package main

import (
	"auth_api/controller"
	"auth_api/db"
	"auth_api/repository"
	"auth_api/usecase"
	"bytes"
	"encoding/json"
	"io"
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

		responseBody, err := io.ReadAll(rr.Result().Body)
		if err != nil {
			panic(err)
		}
		responseBodyMap := make(map[string]interface{})
		json.NewDecoder(bytes.NewReader(responseBody)).Decode(&responseBodyMap)
		if rr.Code == http.StatusOK {
			if _, ok := responseBodyMap["token"]; !ok {
				t.Errorf("Expected token to be in response body")
			}
		}

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
	bodyJSON, _ := json.Marshal(
		map[string]string{
			"email":    "teste@gmail.com",
			"password": "123",
		},
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyJSON))

	userController.Login(rr, req)

	responseBody, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	responseBodyMap := make(map[string]interface{})
	json.NewDecoder(bytes.NewReader(responseBody)).Decode(&responseBodyMap)
	token, ok := responseBodyMap["token"].(string)
	if !ok {
		t.Errorf("Expected token to be in response body")
	}

	tests := []TestCase{
		{
			map[string]string{
				"content": "test content",
			},
			http.StatusCreated,
		},
		{
			map[string]string{
				"content": "",
			},
			http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		// create post
		bodyJSON, _ = json.Marshal(
			test.body,
		)

		rr = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Authorization", token)

		postController.CreatePost(rr, req)

		if rr.Code != test.expectedStatusCode {
			t.Errorf("Expected status code %d, got %d", test.expectedStatusCode, rr.Code)
		}
	}
}
