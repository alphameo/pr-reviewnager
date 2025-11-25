package app

import "github.com/alphameo/pr-reviewnager/internal/domain"

type TeamDTO struct {
	ID      domain.ID
	Name    string
	UserIDs []domain.ID
}

type TeamWithUsersDTO struct {
	TeamName  string
	TeamUsers []*UserDTO
}
