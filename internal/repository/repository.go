package repository

import (
	"JWT_auth/internal/model"
	"context"
	"errors"
)

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

//save user in db
func (r *Repository) SaveUser(user *model.User) (string, error) {
	var id string
	q := `INSERT INTO users(username,email,phone,password)
    VALUES($1,$2,$3,$4)
	RETURNING id;`
	r.db.QueryRow(context.Background(), q, user.Username, user.Email, user.Phone, user.Password).Scan(&id)
	if len(id) == 0 {
		return "", errors.New("error: The user is not saved")
	}
	return id, nil
}

//get user from db
func (r *Repository) GetUser(user *model.User) (string, string, error) {
	var id, role string
	q := `SELECT id,role FROM users
	WHERE
		email=$1 AND password=$2;`
	r.db.QueryRow(context.Background(), q, user.Email, user.Password).Scan(&id, &role)
	if id == "" {
		return "", "", errors.New("error: user not found")
	}
	return id, role, nil
}
