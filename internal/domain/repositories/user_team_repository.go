package repositories

import (
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserTeamRepository interface {
	CreateTeamAndModifyUsers(team *e.Team, users []e.User) error
	FindTeamByTeammateID(userID v.ID) (*e.Team, error)
	FindActiveUsersByTeamID(teamID v.ID) ([]e.User, error)
}
