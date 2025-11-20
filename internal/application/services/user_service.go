// Package services provides application services
package services

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserService interface {
	RegisterUser(user *dto.UserDTO) error
	UnregisterUserByID(userID v.ID) error
	ListUsers() ([]*dto.UserDTO, error)
}

type DefaultUserService struct {
	userRepo r.UserRepository
}

func NewDefaultUserService(userRepository r.UserRepository) (*DefaultUserService, error) {
	if userRepository == nil {
		return nil, errors.New("userRepository cannot be nil")
	}

	return &DefaultUserService{userRepo: userRepository}, nil
}

func (s *DefaultUserService) RegisterUser(user *dto.UserDTO) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	entity, err := mappers.UserToEntity(user)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(entity)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUserService) UnregisterUserByID(userID v.ID) error {
	err := s.userRepo.DeleteByID(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultUserService) ListUsers() ([]*dto.UserDTO, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return mappers.UsersToDTOs(users)
}
