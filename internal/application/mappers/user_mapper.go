package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func UserToDTO(user *e.User) (*dto.UserDTO, error) {
	if user == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &dto.UserDTO{
		ID:     user.ID(),
		Name:   user.Name(),
		Active: user.Active(),
	}, nil
}

func UserToEntity(dto *dto.UserDTO) (*e.User, error) {
	return e.NewExistingUser(dto)
}

func UsersToDTOs(users []*e.User) ([]*dto.UserDTO, error) {
	return EntitiesToDTOs(users, UserToDTO)
}

func UsersToEntities(dtos []*dto.UserDTO) ([]*e.User, error) {
	return DTOsToEntities(dtos, UserToEntity)
}
