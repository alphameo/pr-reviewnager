package services

import (
	"github.com/alphameo/pr-reviewnager/internal/application/dto"
)

type TeamService interface {
	CreateTeam(team dto.TeamDTO) error
	FindTeamByName(name string) (dto.TeamDTO, error)
}
