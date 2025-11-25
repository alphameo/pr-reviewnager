package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func TeamToDTO(entity *e.Team) (*dto.Team, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &dto.Team{
		ID:      entity.ID(),
		Name:    entity.Name(),
		UserIDs: entity.UserIDs(),
	}, nil
}

func TeamToEntity(dto *dto.Team) (*e.Team, error) {
	return e.NewExistingTeam(dto)
}

func TeamsToDTOs(entities []*e.Team) ([]*dto.Team, error) {
	return EntitiesToDTOs(entities, TeamToDTO)
}

func TeamsToEntities(dtos []*dto.Team) ([]*e.Team, error) {
	return DTOsToEntities(dtos, TeamToEntity)
}
