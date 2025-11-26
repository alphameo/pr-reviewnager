package app

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

type UserService interface {
	RegisterUser(user *NewUserDTO) error
	UnregisterUserByID(userID domain.ID) error
	ListUsers() ([]*UserDTO, error)
}

type DefaultUserService struct {
	userRepo domain.UserRepository
}

func NewDefaultUserService(userRepository domain.UserRepository) (*DefaultUserService, error) {
	if userRepository == nil {
		return nil, errors.New("userRepository cannot be nil")
	}

	return &DefaultUserService{userRepo: userRepository}, nil
}

func (s *DefaultUserService) RegisterUser(user *NewUserDTO) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	entity, err := domain.NewUser(user.Name, user.Active)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(entity)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUserService) UnregisterUserByID(userID domain.ID) error {
	err := s.userRepo.DeleteByID(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUserService) ListUsers() ([]*UserDTO, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return UsersToDTOs(users)
}
