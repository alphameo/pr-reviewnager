package dto

import v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"

type Team struct {
	ID      v.ID
	Name    string
	UserIDs []v.ID
}
