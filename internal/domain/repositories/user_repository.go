package repositories

import (
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserRepository interface {
	Repository[e.User, v.ID]
	FindActiveUsersByTeamID(teamID v.ID) ([]e.User, error)
}
