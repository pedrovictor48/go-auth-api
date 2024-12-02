package controller

import (
	"auth_api/model"
	"auth_api/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) UserController {
	return UserController{client}
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")

	var user UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}
	userRepository := repository.NewUserRepository(c.client)
	existingUser, err := userRepository.GetUserByEmail(user.Email)
	if err == repository.ErrUserNotFound {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, fmt.Sprintf("Senha incorreta, %v", err), http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": existingUser.ID,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {

	var user model.UserRegister
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}
	userRepository := repository.NewUserRepository(c.client)
	err = userRepository.CreateUser(user)
	if err == repository.ErrEmailAlreadyExists {
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar usuário: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
