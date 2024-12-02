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
)

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return UserRepository{client}
}

func (r *UserRepository) CreateUser(user model.User) error {
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
		return err
	}

	user.Password = string(hash)
	_, err = collection.InsertOne(context.TODO(), bson.M{
		"email":     user.Email,
		"password":  user.Password,
		"name":      user.Name,
		"birthdate": user.Birthdate,
	})
	if err != nil {
		return err
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
	panic(err)
	return model.User{}, ErrUserNotFound
}
