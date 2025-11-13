package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func UserToDTO(entity *e.User) *dto.UserDTO {
	dto := dto.UserDTO{
		ID:     entity.ID(),
		Name:   entity.Name(),
		Active: entity.Active(),
	}

	return &dto
}

func UserDTOToEntity(dto *dto.UserDTO) (*e.User, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}
	entity := e.NewUserWithID(dto.ID, dto.Name, dto.Active)

	return entity, nil
}
