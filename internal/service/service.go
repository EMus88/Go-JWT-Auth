package service

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	SaveUser(user *model.User) (string, error)
}

type Service struct {
	Repository
}

func NewService(r *repository.Repository) *Service {
	return &Service{Repository: r}
}

func (s *Service) SignIn(user *model.User) error {

	//hashing the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)

	//try savign user in DB
	id, err := s.Repository.SaveUser(user)
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
