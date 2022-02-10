package service

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/repository"
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

type Auth struct {
	Repository
}

func NewAuth(repos *repository.Repository) *Auth {
	return &Auth{Repository: repos}
}

func (a *Auth) CreateUser(user *model.User) error {
	//hashing the password
	user.Password = a.HashingPassword(user.Password)

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
	user.Password = "******"

	return nil
}

func (a *Auth) GenerateTokenPair(id string) (string, string, error) {
	//create access token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		Subject:   "access",
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", "", err
	}
	//create refresh token
	rtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Hour * 1000).Unix(),
		Subject:   "refresh",
	})
	rToken, err := rtClaims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", "", err
	}
	return token, rToken, nil
}

func (a *Auth) HashingPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	hash := h.Sum([]byte(os.Getenv("SALT")))
	return fmt.Sprintf("%x", hash)
}
