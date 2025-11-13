package repositories

import (
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type TeamRepository interface {
	Repository[e.Team, v.ID]
	FindTeamByName(name string) (e.Team, error)
	FindTeamByTeammateID(teamateID v.ID) (e.Team, error)
}
