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

func UsersToDTOs(entities []e.User) []dto.UserDTO {
	dtos := make([]dto.UserDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *UserToDTO(&entity)
	}

	return dtos
}

func UserDTOToEntity(dto *dto.UserDTO) (*e.User, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	entity := e.NewUserWithID(
		dto.ID,
		dto.Name,
		dto.Active,
	)

	return entity, nil
}

func UserDTOsToEntities(dtos []dto.UserDTO) ([]e.User, error) {
	entities := make([]e.User, len(dtos))
	for i, dto := range dtos {
		entity, err := UserDTOToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = *entity
	}

	return entities, nil
}
