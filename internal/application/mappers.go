package app

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

func EntitiesToDTOs[ENTITY any, DTO any](entities []*ENTITY, mapFunc func(*ENTITY) (*DTO, error)) ([]*DTO, error) {
	dtos := make([]*DTO, len(entities))
	for i, entity := range entities {
		dto, err := mapFunc(entity)
		if err != nil {
			return nil, err
		}
		dtos[i] = dto
	}

	return dtos, nil
}

func DTOsToEntities[ENTITY any, DTO any](dtos []*DTO, mapFunc func(*DTO) (*ENTITY, error)) ([]*ENTITY, error) {
	entities := make([]*ENTITY, len(dtos))
	for i, dto := range dtos {
		entity, err := mapFunc(dto)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func PullRequestToDTO(entity *domain.PullRequest) (*domain.PullRequestDTO, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &domain.PullRequestDTO{
		ID:          entity.ID(),
		Title:       entity.Title(),
		AuthorID:    entity.AuthorID(),
		CreatedAt:   entity.CreatedAt(),
		Status:      entity.Status().String(),
		MergedAt:    entity.MergedAt(),
		ReviewerIDs: entity.ReviewerIDs(),
	}, nil
}

func PullRequestToEntity(dto *domain.PullRequestDTO) (*domain.PullRequest, error) {
	return domain.NewExistingPullRequest(dto)
}

func PullRequestsToDTOs(entities []*domain.PullRequest) ([]*domain.PullRequestDTO, error) {
	return EntitiesToDTOs(entities, PullRequestToDTO)
}

func PullRequestsToEntities(dtos []*domain.PullRequestDTO) ([]*domain.PullRequest, error) {
	return DTOsToEntities(dtos, PullRequestToEntity)
}

func TeamToDTO(entity *domain.Team) (*domain.TeamDTO, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &domain.TeamDTO{
		ID:      entity.ID(),
		Name:    entity.Name(),
		UserIDs: entity.UserIDs(),
	}, nil
}

func TeamToEntity(dto *domain.TeamDTO) (*domain.Team, error) {
	return domain.NewExistingTeam(dto)
}

func TeamsToDTOs(entities []*domain.Team) ([]*domain.TeamDTO, error) {
	return EntitiesToDTOs(entities, TeamToDTO)
}

func TeamsToEntities(dtos []*domain.TeamDTO) ([]*domain.Team, error) {
	return DTOsToEntities(dtos, TeamToEntity)
}

func UserToDTO(user *domain.User) (*domain.UserDTO, error) {
	if user == nil {
		return nil, errors.New("entity cannot be nil")
	}

	return &domain.UserDTO{
		ID:     user.ID(),
		Name:   user.Name(),
		Active: user.Active(),
	}, nil
}

func UserToEntity(dto *domain.UserDTO) (*domain.User, error) {
	return domain.NewExistingUser(dto)
}

func UsersToDTOs(users []*domain.User) ([]*domain.UserDTO, error) {
	return EntitiesToDTOs(users, UserToDTO)
}

func UsersToEntities(dtos []*domain.UserDTO) ([]*domain.User, error) {
	return DTOsToEntities(dtos, UserToEntity)
}
