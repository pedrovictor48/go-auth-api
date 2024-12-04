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
