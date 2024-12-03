package controller

import (
	"auth_api/repository"
	"auth_api/usecase"
	"encoding/json"
	"net/http"
)

type UserController struct {
	usecase usecase.UserUsecase
}

func NewUserController(usecase usecase.UserUsecase) UserController {
	return UserController{usecase}
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var user repository.UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	tokenString, err := c.usecase.LoginUser(user)
	if err == repository.ErrUserNotFound {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	} else if err == usecase.ErrEncriptPassword {
		http.Error(w, "Senha inválida", http.StatusUnauthorized)
		return
	} else if err == usecase.ErrGenerateToken {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {

	var user repository.UserRegister
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	err = c.usecase.RegisterUser(user)
	if err == repository.ErrEmailAlreadyExists {
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	} else if err == repository.ErrEncriptPassword {
		http.Error(w, "Erro ao encriptar senha", http.StatusInternalServerError)
		return
	} else if err == repository.ErrInsertUser {
		http.Error(w, "Erro ao inserir usuário", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
