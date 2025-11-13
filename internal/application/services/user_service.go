// Package services provides application services
package services

import (
	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserService interface {
	RegisterUser(user dto.UserDTO) error
	UnregisterUserByID(userID v.ID) error
	ListUsers() ([]dto.UserDTO, error)
	SetUserActiveByID(userID v.ID, active bool) error
}
