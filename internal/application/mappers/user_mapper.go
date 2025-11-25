package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func UserToDTO(user *e.User) (*dto.User, error) {
	if user == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &dto.User{
		ID:     user.ID(),
		Name:   user.Name(),
		Active: user.Active(),
	}, nil
}

func UserToEntity(dto *dto.User) (*e.User, error) {
	return e.NewExistingUser(dto)
}

func UsersToDTOs(users []*e.User) ([]*dto.User, error) {
	return EntitiesToDTOs(users, UserToDTO)
}

func UsersToEntities(dtos []*dto.User) ([]*e.User, error) {
	return DTOsToEntities(dtos, UserToEntity)
}
