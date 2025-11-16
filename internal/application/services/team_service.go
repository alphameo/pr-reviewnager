package services

import (
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	"github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeamWithUsers(teamDTO dto.TeamWithUsersDTO) error
	FindTeamByName(name string) (*dto.TeamWithUsersDTO, error)
	SetUserActiveByID(userID v.ID, active bool) (*dto.UserWithTeamNameDTO, error)
}

var (
	ErrTeamExists      error = errors.New("team already exists")
	ErrPRExists        error = errors.New("pull request already exists")
	ErrPRAlreadyMerged error = errors.New("pull request already merged")
	ErrNoCandidate     error = errors.New("no active candidate for assigning")
	ErrNotFound        error = errors.New("resource not found")
	ErrNotAssigned     error = errors.New("reviewer is not assigned to PR")
)

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

func (s *DefaultTeamService) CreateTeamWithUsers(teamDTO dto.TeamWithUsersDTO) error {
	existingTeam, err := s.teamRepo.FindByName(teamDTO.TeamName)
	if err != nil {
		return ErrTeamExists
	}
	if existingTeam != nil {
		return fmt.Errorf("team with name=%s already exists", teamDTO.TeamName)
	}

	var team *entities.Team
	if teamDTO.ID == v.ID(uuid.Nil) {
		team = entities.NewTeam(teamDTO.TeamName)
	} else {
		team = entities.NewTeamWithID(teamDTO.ID, teamDTO.TeamName)
	}
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

func (s *DefaultTeamService) FindTeamByName(name string) (*dto.TeamWithUsersDTO, error) {
	team, err := s.teamRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, ErrNotFound
	}
	users := make([]*dto.UserDTO, len(team.UserIDs()))
	for i, userID := range team.UserIDs() {
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, err
		}
		dto, err := mappers.UserToDTO(user)
		if err != nil {
			return nil, err
		}
		users[i] = dto
	}

	res := dto.TeamWithUsersDTO{
		ID:        team.ID(),
		TeamName:  team.Name(),
		TeamUsers: users,
	}
	return &res, nil
}

func (s *DefaultTeamService) SetUserActiveByID(userID v.ID, active bool) (*dto.UserWithTeamNameDTO, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("%w: no such user with id=%d", ErrNotFound, userID)
	}
	user.SetActive(active)

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	team, err := s.teamRepo.FindTeamByTeammateID(userID)
	if err != nil {
		return nil, err
	}
	userDTO, err := mappers.UserToDTO(user)
	if err != nil {
		return nil, err
	}
	dto := dto.UserWithTeamNameDTO{
		User:     userDTO,
		TeamName: team.Name(),
	}
	return &dto, nil
}
