package services

import (
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeamWithUsers(teamDTO *TeamWithUsersDTO) error
	FindTeamByName(name string) (*TeamWithUsersDTO, error)
	SetUserActiveByID(userID v.ID, active bool) (*UserWithTeamNameDTO, error)
}

type TeamWithUsersDTO struct {
	ID        v.ID
	TeamName  string
	TeamUsers []*dto.UserDTO
}

type UserWithTeamNameDTO struct {
	User     *dto.UserDTO
	TeamName string
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

	return &DefaultTeamService{
		teamRepo: teamRepository,
		userRepo: userRepository,
	}, nil
}

func (s *DefaultTeamService) CreateTeamWithUsers(teamDTO *TeamWithUsersDTO) error {
	if teamDTO == nil {
		return errors.New("dto cannot be nil")
	}
	existingTeam, err := s.teamRepo.FindByName(teamDTO.TeamName)
	if err != nil {
		return ErrTeamExists
	}
	if existingTeam != nil {
		return fmt.Errorf("team with name=%s already exists", teamDTO.TeamName)
	}

	var team *e.Team
	if teamDTO.ID == v.ID(uuid.Nil) {
		team, err = e.NewTeam(teamDTO.TeamName)
	} else {
		tDTO := dto.TeamDTO{ID: teamDTO.ID, Name: teamDTO.TeamName, UserIDs: nil}
		team, err = e.NewExistingTeam(&tDTO)
	}
	if err != nil {
		return err
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

func (s *DefaultTeamService) FindTeamByName(name string) (*TeamWithUsersDTO, error) {
	team, err := s.teamRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, ErrNotFound
	}
	users := make([]*dto.UserDTO, len(team.UserIDs))
	for i, userID := range team.UserIDs {
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return &TeamWithUsersDTO{
		ID:        team.ID,
		TeamName:  team.Name,
		TeamUsers: users,
	}, nil
}

func (s *DefaultTeamService) SetUserActiveByID(userID v.ID, active bool) (*UserWithTeamNameDTO, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, fmt.Errorf("%w: no such user with id=%d", ErrNotFound, userID)
	}

	user, err := e.NewExistingUser(u)
	if err != nil {
		return nil, err
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

	return &UserWithTeamNameDTO{
		User:     userDTO,
		TeamName: team.Name,
	}, nil
}
