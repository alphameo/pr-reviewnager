package repositories

import (
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type TeamRepository interface {
	Repository[e.Team, dto.Team, v.ID]
	FindByName(teamName string) (*dto.Team, error)
	CreateTeamAndModifyUsers(team *e.Team, users []*e.User) error
	FindTeamByTeammateID(userID v.ID) (*dto.Team, error)
	FindActiveUsersByTeamID(teamID v.ID) ([]*dto.User, error)
}
