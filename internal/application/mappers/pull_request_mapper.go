// Package mappers provide mapping between domain entitites, value objects and application dto
package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

func PullRequestToDTO(entity *e.PullRequest) *dto.PullRequestDTO {
	status := entity.Status()
	dto := dto.PullRequestDTO{
		ID:          entity.ID(),
		Title:       entity.Title(),
		AuthorID:    entity.AuthorID(),
		Status:      status.String(),
		MergedAt:    entity.MergedAt(),
		ReviewerIDs: entity.ReviewerIDs(),
	}

	return &dto
}

func PullRequestsToDTOs(entities []e.PullRequest) ([]dto.PullRequestDTO, error) {
	if entities == nil {
		return nil, errors.New("entitites cannot be nil")
	}

	dtos := make([]dto.PullRequestDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *PullRequestToDTO(&entity)
	}

	return dtos, nil
}

func PullRequestDTOToEntity(dto *dto.PullRequestDTO) (*e.PullRequest, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	status, err := v.NewPRStatusFromString(dto.Status)
	if err != nil {
		return nil, err
	}
	entity := e.NewPullRequestWithID(dto.ID, dto.Title, dto.AuthorID, status, dto.MergedAt)
	for _, id := range dto.ReviewerIDs {
		entity.AssignReviewer(id)
	}

	return entity, nil
}

func PullRequestDTOsToEntities(dtos []dto.PullRequestDTO) ([]e.PullRequest, error) {
	if dtos == nil {
		return nil, errors.New("dtos cannot be nil")
	}

	entities := make([]e.PullRequest, len(dtos))
	for i, dto := range dtos {
		entity, err := PullRequestDTOToEntity(&dto)
		if err != nil {
			return nil, err
		}
		entities[i] = *entity
	}

	return entities, nil
}
