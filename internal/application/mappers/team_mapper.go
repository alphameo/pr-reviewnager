package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func TeamToDTO(entity *e.Team) (*dto.TeamDTO, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &dto.TeamDTO{
		ID:      entity.ID(),
		Name:    entity.Name(),
		UserIDs: entity.UserIDs(),
	}, nil
}

func TeamToEntity(dto *dto.TeamDTO) (*e.Team, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	return e.NewExistingTeam(dto.ID, dto.Name, dto.UserIDs)
}

func TeamsToDTOs(entities []*e.Team) ([]*dto.TeamDTO, error) {
	return EntitiesToDTOs(entities, TeamToDTO)
}

func TeamsToEntities(dtos []*dto.TeamDTO) ([]*e.Team, error) {
	return DTOsToEntities(dtos, TeamToEntity)
}
