package usecase

import (
	"auth_api/repository"
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEncriptPassword = errors.New("erro ao encriptar senha")
	ErrGenerateToken   = errors.New("erro ao gerar token")
)

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return UserUsecase{
		repository: repo,
	}
}

func (u *UserUsecase) LoginUser(user repository.UserLogin) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	userRepository := u.repository
	existingUser, err := userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return "", ErrEncriptPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": existingUser.ID,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", ErrGenerateToken
	}
	return tokenString, nil
}

func (u *UserUsecase) RegisterUser(user repository.UserRegister) error {
	var err error
	userRepository := u.repository
	err = userRepository.CreateUser(user)

	return err
}
