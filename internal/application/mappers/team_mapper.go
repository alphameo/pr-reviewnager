package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
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
	return e.NewExistingTeam(dto)
}

func TeamsToDTOs(entities []*e.Team) ([]*dto.TeamDTO, error) {
	return EntitiesToDTOs(entities, TeamToDTO)
}

func TeamsToEntities(dtos []*dto.TeamDTO) ([]*e.Team, error) {
	return DTOsToEntities(dtos, TeamToEntity)
}
