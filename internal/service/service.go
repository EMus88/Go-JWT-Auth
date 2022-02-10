package service

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/repository"
)

type Repository interface {
	//auth methods
	SaveUser(user *model.User) (string, error)
	GetUser(user *model.User) (string, error)
}

type Service struct {
	Repository
	Auth
}

func NewService(r *repository.Repository) *Service {
	return &Service{Repository: r, Auth: *NewAuth(r)}
}
