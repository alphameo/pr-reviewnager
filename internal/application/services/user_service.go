// Package services provides application services
package services

import (
	"errors"
	"fmt"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserService interface {
	RegisterUser(user dto.UserDTO) error
	UnregisterUserByID(userID v.ID) error
	ListUsers() ([]dto.UserDTO, error)
	SetUserActiveByID(userID v.ID, active bool) error
}

type DefaultUserService struct {
	userRepo r.UserRepository
}

func NewDefaulUserService(userRepository *r.UserRepository) (*DefaultUserService, error) {
	if userRepository == nil {
		return nil, errors.New("userRepository cannot be nil")
	}
	s := DefaultUserService{userRepo: *userRepository}
	return &s, nil
}

func (s *DefaultUserService) RegisterUser(user *dto.UserDTO) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	entity, err := mappers.UserDTOToEntity(user)
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

func (s *DefaultUserService) ListUsers() ([]dto.UserDTO, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return mappers.UsersToDTOs(users), nil
}

func (s *DefaultUserService) SetUserActiveByID(userID v.ID, active bool) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("no such user with id=%d", userID)
	}
	user.SetActive(active)

	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}
