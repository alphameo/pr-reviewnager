// Package mappers provide mapping between domain entitites, value objects and application dto
package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
)

func PullRequestToDTO(entity *e.PullRequest) (*dto.PullRequestDTO, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &dto.PullRequestDTO{
		ID:          entity.ID(),
		Title:       entity.Title(),
		AuthorID:    entity.AuthorID(),
		CreatedAt:   entity.CreatedAt(),
		Status:      entity.Status().String(),
		MergedAt:    entity.MergedAt(),
		ReviewerIDs: entity.ReviewerIDs(),
	}, nil
}

func PullRequestToEntity(dto *dto.PullRequestDTO) (*e.PullRequest, error) {
	return e.NewExistingPullRequest(dto)
}

func PullRequestsToDTOs(entities []*e.PullRequest) ([]*dto.PullRequestDTO, error) {
	return EntitiesToDTOs(entities, PullRequestToDTO)
}

func PullRequestsToEntities(dtos []*dto.PullRequestDTO) ([]*e.PullRequest, error) {
	return DTOsToEntities(dtos, PullRequestToEntity)
}
