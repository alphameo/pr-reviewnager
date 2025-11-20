// Package mappers provide mapping between domain entitites, value objects and application dto
package mappers

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
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
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	status, err := v.NewPRStatusFromString(dto.Status)
	if err != nil {
		return nil, err
	}
	return e.NewExistingPullRequest(
		dto.ID,
		dto.Title,
		dto.AuthorID,
		dto.CreatedAt,
		status,
		nil,
		dto.ReviewerIDs,
	)
}

func PullRequestsToDTOs(entities []*e.PullRequest) ([]*dto.PullRequestDTO, error) {
	return EntitiesToDTOs(entities, PullRequestToDTO)
}

func PullRequestsToEntities(dtos []*dto.PullRequestDTO) ([]*e.PullRequest, error) {
	return DTOsToEntities(dtos, PullRequestToEntity)
}
