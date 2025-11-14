package services

import (
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	"github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
)

type TeamService interface {
	CreateTeamWithUsers(teamDTO dto.CreateTeamWithUsersDTO) error
	FindTeamByName(name string) (*dto.TeamDTO, error)
}

type DefaultTeamService struct {
	teamRepo     r.TeamRepository
	userRepo     r.UserRepository
	userTeamRepo r.UserTeamRepository
}

func (s *DefaultTeamService) CreateTeamWithUsers(teamDTO dto.CreateTeamWithUsersDTO) error {
	existingTeam, err := s.teamRepo.FindTeamByName(teamDTO.TeamName)
	if err != nil {
		return err
	}
	if existingTeam != nil {
		return fmt.Errorf("team with name=%s already exists", teamDTO.TeamName)
	}

	team := entities.NewTeam(teamDTO.TeamName)
	users, err := mappers.UserDTOsToEntities(teamDTO.TeamUsers)
	if err != nil {
		return err
	}

	for _, user := range users {
		team.AddUser(user.ID())
	}

	s.userTeamRepo.CreateTeamAndModifyUsers(team, users)
	return nil
}

func (s *DefaultTeamService) FindTeamByName(name string) (*dto.TeamDTO, error) {
	team, err := s.teamRepo.FindTeamByName(name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, nil
	}
	dto := mappers.TeamToDTO(team)
	return dto, nil
}
