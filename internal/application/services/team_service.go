package services

import (
	"github.com/alphameo/pr-reviewnager/internal/application/dto"
)

type TeamService interface {
	CreateTeamWithUsers(teamDTO dto.CreateTeamWithUsersDTO) error
	FindTeamByName(name string) (*dto.TeamDTO, error)
}
