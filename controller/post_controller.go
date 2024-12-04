package controller

import (
	"auth_api/model"
	"auth_api/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PostController struct {
	postUsecase usecase.PostUsecase
}

func NewPostController(postUsecase usecase.PostUsecase) PostController {
	return PostController{
		postUsecase: postUsecase,
	}
}

type CreatePostRequest struct {
	Content string `json:"content"`
}

func (p *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var postRequest CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&postRequest)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	if postRequest.Content == "" {
		http.Error(w, "Conteúdo do post não pode ser vazio", http.StatusBadRequest)
		return
	}

	//validate authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Token não encontrado", http.StatusUnauthorized)
		return
	}
	fmt.Println(tokenString)

	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Token inválido, %v", err), http.StatusUnauthorized)
		return
	}
	var authorId string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		authorId = claims["id"].(string)
	} else {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	newPost := model.Post{
		Content: postRequest.Content,
		Author:  authorId,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	}

	err = p.postUsecase.CreatePost(newPost)
	if err != nil {
		http.Error(w, "Erro ao criar post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

type GetPostsRequest struct {
	AuthorId string `json:"author_id"`
}

func (p *PostController) GetPostsById(w http.ResponseWriter, r *http.Request) {
	var getPostsRequest GetPostsRequest
	err := json.NewDecoder(r.Body).Decode(&getPostsRequest)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}
	posts, err := p.postUsecase.GetPostsById(getPostsRequest.AuthorId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar posts, %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
