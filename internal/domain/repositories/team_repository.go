package repositories

import (
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type TeamRepository interface {
	Repository[e.Team, dto.TeamDTO, v.ID]
	FindByName(teamName string) (*dto.TeamDTO, error)
	CreateTeamAndModifyUsers(team *e.Team, users []*e.User) error
	FindTeamByTeammateID(userID v.ID) (*dto.TeamDTO, error)
	FindActiveUsersByTeamID(teamID v.ID) ([]*dto.UserDTO, error)
}
