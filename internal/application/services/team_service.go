package services

import (
	"errors"
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
	teamRepo r.TeamRepository
	userRepo r.UserRepository
}

func NewDefaultTeamService(
	teamRepository r.TeamRepository,
	userRepository r.UserRepository,
) (*DefaultTeamService, error) {
	if teamRepository == nil {
		return nil, errors.New("teamRepository cannot be nil")
	}
	if userRepository == nil {
		return nil, errors.New("userRepository cannot be nil")
	}

	s := DefaultTeamService{
		teamRepo: teamRepository,
		userRepo: userRepository,
	}
	return &s, nil
}

func (s *DefaultTeamService) CreateTeamWithUsers(teamDTO dto.CreateTeamWithUsersDTO) error {
	existingTeam, err := s.teamRepo.FindByName(teamDTO.TeamName)
	if err != nil {
		return err
	}
	if existingTeam != nil {
		return fmt.Errorf("team with name=%s already exists", teamDTO.TeamName)
	}

	team := entities.NewTeam(teamDTO.TeamName)
	users, err := mappers.UsersToEntities(teamDTO.TeamUsers)
	if err != nil {
		return err
	}

	for _, user := range users {
		team.AddUser(user.ID())
	}

	s.teamRepo.CreateTeamAndModifyUsers(team, users)
	return nil
}

func (s *DefaultTeamService) FindTeamByName(name string) (*dto.TeamDTO, error) {
	team, err := s.teamRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, nil
	}
	dto, err := mappers.TeamToDTO(team)
	if err != nil {
		return nil, err
	}
	return dto, nil
}
