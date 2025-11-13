package dto

import v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"

type TeamDTO struct {
	ID      v.ID
	Name    string
	UserIDs []v.ID
}
