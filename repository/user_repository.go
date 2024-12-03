package repository

import (
	"auth_api/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrEncriptPassword    = errors.New("erro ao encriptar senha")
	ErrInsertUser         = errors.New("erro ao inserir usuário")
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

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return UserRepository{client}
}

func (r *UserRepository) CreateUser(user UserRegister) error {
	var err error
	collection := r.client.Database("testdb").Collection("users")
	//
	var existingUser model.User
	err = collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrEncriptPassword
	}

	user.Password = string(hash)
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return ErrInsertUser
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (model.User, error) {
	var err error
	collection := r.client.Database("testdb").Collection("users")
	var existingUser model.User
	err = collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return existingUser, nil
	}
	return model.User{}, ErrUserNotFound
}
