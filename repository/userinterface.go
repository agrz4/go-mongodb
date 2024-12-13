package repository

import "go-mongodb/model"

type UserInterface interface {
	CreateUser(model.User) (string, error)
	GetUserByID(string) (model.User, error)
	GetAllUsers() ([]model.User, error)
	UpdateUserAgeByID(string, int) (int, error)
	DeleteUserByID(string) (int, error)
	DeleteAllUsers() (int, error)
}
