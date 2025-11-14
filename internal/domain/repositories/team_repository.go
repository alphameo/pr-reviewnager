package repositories

import (
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type TeamRepository interface {
	Repository[e.Team, v.ID]
	FindTeamByName(teamName string) (*e.Team, error)
}
