package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Name      string `json:"name" bson:"name"`
	Birthdate string `json:"birthdate" bson:"birthdate"`
	Gender    string `json:"gender" bson:"gender"`
}

type User struct {
	Email     string   `json:"email" bson:"email"`
	Password  string   `json:"password" bson:"password"`
	Name      string   `json:"name" bson:"name"`
	Birthdate string   `json:"birthdate" bson:"birthdate"`
	Gender    string   `json:"gender" bson:"gender"`
	ID        string   `json:"id" bson:"_id"`
	Friends   []string `json:"friends" bson:"friends"`
}

var client *mongo.Client

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	secret := os.Getenv("JWT_SECRET")

	// Set client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/friend", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodGet == r.Method {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			token, _, err := new(jwt.Parser).ParseUnverified(authHeader, jwt.MapClaims{})
			if err != nil {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(token.Claims)
			userID := token.Claims.(jwt.MapClaims)["id"].(string)
			objectID, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				log.Fatalf("Erro ao converter ID para ObjectID: %v", err)
			}

			collection := client.Database("testdb").Collection("users")
			var user User
			err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
			if err != nil {
				http.Error(w, fmt.Sprintf("Usuário não encontrado, %v", err), http.StatusNotFound)
				return
			}
			friends := user.Friends

			json.NewEncoder(w).Encode(friends)

		} else if http.MethodPost == r.Method {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}
			token, _, err := new(jwt.Parser).ParseUnverified(authHeader, jwt.MapClaims{})
			if err != nil {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(token.Claims)
			userID := token.Claims.(jwt.MapClaims)["id"].(string)
			objectID, err := primitive.ObjectIDFromHex(userID)
			if err != nil {
				log.Fatalf("Erro ao converter ID para ObjectID: %v", err)
			}

			collection := client.Database("testdb").Collection("users")
			var user User
			err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
			if err != nil {
				http.Error(w, fmt.Sprintf("Usuário não encontrado, %v", err), http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			// Decodificar o corpo da requisição JSON
			var user UserLogin
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&user)
			if err != nil {
				http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
				return
			}
			collection := client.Database("testdb").Collection("users")

			var existingUser User
			err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
			if err != nil {
				http.Error(w, "Usuário não encontrado", http.StatusNotFound)
				return
			}
			//json.NewEncoder(w).Encode(existingUser)

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
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			var user UserRegister
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&user)
			if err != nil {
				http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
				return
			}

			collection := client.Database("testdb").Collection("users")
			//
			var existingUser UserRegister
			err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
			if err == nil {
				http.Error(w, "Email já cadastrado", http.StatusConflict)
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Erro ao criptografar senha", http.StatusInternalServerError)
				return
			}

			user.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				http.Error(w, fmt.Sprintf("Erro ao registrar usuário: %v", err), http.StatusInternalServerError)
				return
			}

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
