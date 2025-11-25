package infra

import "github.com/alphameo/pr-reviewnager/internal/domain"

type TeamDTO struct {
	ID      domain.ID
	Name    string
	UserIDs []domain.ID
}
