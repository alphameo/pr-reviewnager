package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func TeamToDTO(entity *e.Team) *dto.TeamDTO {
	dto := dto.TeamDTO{
		ID:      entity.ID(),
		Name:    entity.Name(),
		UserIDs: entity.UserIDs(),
	}

	return &dto
}

func TeamsToDTOs(entities []e.Team) ([]dto.TeamDTO, error) {
	if entities == nil {
		return nil, errors.New("entitites cannot be nil")
	}

	dtos := make([]dto.TeamDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *TeamToDTO(&entity)
	}
	
	return dtos, nil
}

func TeamDTOToEntity(dto *dto.TeamDTO) (*e.Team, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	entity := e.NewTeamWithID(dto.ID, dto.Name)
	for _, id := range dto.UserIDs {
		entity.AddUser(id)
	}

	return entity, nil
}

func TeamDTOsToEntities(dtos []dto.TeamDTO) ([]e.Team, error) {
	if dtos == nil {
		return nil, errors.New("dtos cannot be nil")
	}

	entities := make([]e.Team, len(dtos))
	for i, dto := range dtos {
		entity, err := TeamDTOToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = *entity
	}

	return entities, nil
}
