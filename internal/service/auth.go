package service

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/repository"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Repository
}

func NewAuth(repos *repository.Repository) *Auth {
	return &Auth{Repository: repos}
}

func (a *Auth) CreateUser(user *model.User) error {
	//hashing the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)

	//try saving user in DB
	id, err := a.Repository.SaveUser(user)
	if err != nil {
		return err
	}
	//convert uuid
	uuidID, err := uuid.FromString(id)
	if err != nil {
		return err
	}
	//set id for response
	user.ID = uuidID
	user.Password = ""

	return nil
}

func (a *Auth) GenerateTokenPair(id string) (string, string, error) {
	//create access token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		Subject:   "access",
	})
	token, err := claims.SignedString(os.Getenv("SECRET"))
	if err != nil {
		return "", "", err
	}
	//create refresh token
	rtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Hour * 1000).Unix(),
		Subject:   "refresh",
	})
	rToken, err := rtClaims.SignedString(os.Getenv("SECRET"))
	if err != nil {
		return "", "", err
	}
	return token, rToken, nil
}
