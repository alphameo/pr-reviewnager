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
