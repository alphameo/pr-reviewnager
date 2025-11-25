package app

import "github.com/alphameo/pr-reviewnager/internal/domain"

type UserDTO struct {
	ID     domain.ID
	Name   string
	Active bool
}
