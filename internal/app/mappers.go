package app

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

// Common

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

func DTOsToDomain[ENTITY any, DTO any](dtos []*DTO, mapFunc func(*DTO) (*ENTITY, error)) ([]*ENTITY, error) {
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

var (
	ErrNilDTO       error = errors.New("dto cannot be nil")
	ErrNilDomainObj error = errors.New("domain object cannot be nil")
)

// To DTO

func PullRequestToDTO(entity *domain.PullRequest) (*PullRequestDTO, error) {
	if entity == nil {
		return nil, ErrNilDomainObj
	}

	return &PullRequestDTO{
		ID:          entity.ID(),
		Title:       entity.Title(),
		AuthorID:    entity.AuthorID(),
		CreatedAt:   entity.CreatedAt(),
		Status:      entity.Status().String(),
		MergedAt:    entity.MergedAt(),
		ReviewerIDs: entity.ReviewerIDs(),
	}, nil
}

func PullRequestsToDTOs(entities []*domain.PullRequest) ([]*PullRequestDTO, error) {
	return EntitiesToDTOs(entities, PullRequestToDTO)
}

func TeamToDTO(entity *domain.Team) (*TeamDTO, error) {
	if entity == nil {
		return nil, ErrNilDomainObj
	}

	return &TeamDTO{
		ID:      entity.ID(),
		Name:    entity.Name().Value(),
		UserIDs: entity.UserIDs(),
	}, nil
}

func TeamsToDTOs(entities []*domain.Team) ([]*TeamDTO, error) {
	return EntitiesToDTOs(entities, TeamToDTO)
}

func UserToDTO(user *domain.User) (*UserDTO, error) {
	if user == nil {
		return nil, ErrNilDomainObj
	}

	return &UserDTO{
		ID:     user.ID(),
		Name:   user.Name().Value(),
		Active: user.Active(),
	}, nil
}

func UsersToDTOs(users []*domain.User) ([]*UserDTO, error) {
	return EntitiesToDTOs(users, UserToDTO)
}

// To Domain

func PullRequestToDomain(dto *PullRequestDTO) (*domain.PullRequest, error) {
	if dto == nil {
		return nil, ErrNilDTO
	}

	pr := domain.ExistingPullRequest(
		dto.ID,
		dto.Title,
		dto.AuthorID,
		dto.CreatedAt,
		domain.ExistingPRStatus(dto.Status),
		dto.MergedAt,
		dto.ReviewerIDs,
	)
	if err := pr.Validate(); err != nil {
		return nil, err
	}

	return pr, nil
}

func PullRequestsToDomain(dtos []*PullRequestDTO) ([]*domain.PullRequest, error) {
	return DTOsToDomain(dtos, PullRequestToDomain)
}

func TeamToDomain(dto *TeamDTO) (*domain.Team, error) {
	if dto == nil {
		return nil, ErrNilDTO
	}

	name, err := domain.NewTeamName(dto.Name)
	if err != nil {
		return nil, err
	}

	team := domain.ExistingTeam(dto.ID, name, dto.UserIDs)
	if err := team.Validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func TeamsToEntities(dtos []*TeamDTO) ([]*domain.Team, error) {
	return DTOsToDomain(dtos, TeamToDomain)
}

func UserToDomain(dto *UserDTO) (*domain.User, error) {
	if dto == nil {
		return nil, ErrNilDTO
	}

	name, err := domain.NewUserName(dto.Name)
	if err != nil {
		return nil, err
	}

	user := domain.ExistingUser(dto.ID, name, dto.Active)
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

func UsersToDomain(dtos []*UserDTO) ([]*domain.User, error) {
	return DTOsToDomain(dtos, UserToDomain)
}
