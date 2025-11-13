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
		ReviewerIDs: entity.ReviewerIDs(),
	}

	return &dto
}

func PullRequestDTOToEntity(dto *dto.PullRequestDTO) (*e.PullRequest, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}
	status, err := v.NewPRStatusFromString(dto.Status)
	if err != nil {
		return nil, err
	}
	entity := e.NewPullRequestWithID(dto.ID, dto.Title, dto.AuthorID, status)
	for _, id := range dto.ReviewerIDs {
		entity.AssignReviewer(id)
	}

	return entity, nil
}
