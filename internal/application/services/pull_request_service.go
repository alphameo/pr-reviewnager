package services

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/dto"
	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	s "github.com/alphameo/pr-reviewnager/internal/domain/services"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestService interface {
	CreatePullRequest(pullRequest *dto.PullRequestDTO) (*dto.PullRequestDTO, error)
	MarkAsMerged(pullRequestID v.ID) (*dto.PullRequestDTO, error)
	ReassignReviewer(userID v.ID, pullRequestID v.ID) (*dto.PullRequestWithNewReviewerIDDTO, error)
	FindPullRequestsByReviewer(userID v.ID) ([]*dto.PullRequestDTO, error)
}

type DefaultPullRequestService struct {
	prDomainServ s.PullRequestDomainService
	prRepo       r.PullRequestRepository
}

func NewDefaultPullRequestService(
	pullRequestDomainService s.PullRequestDomainService,
	pullRequestRepository r.PullRequestRepository,
) (*DefaultPullRequestService, error) {
	if pullRequestDomainService == nil {
		return nil, errors.New("pullRequestDomainService cannot bi nil")
	}
	if pullRequestRepository == nil {
		return nil, errors.New("PullRequestRepository cannot be nil")
	}

	s := DefaultPullRequestService{
		prDomainServ: pullRequestDomainService,
		prRepo:       pullRequestRepository,
	}
	return &s, nil
}

func (s *DefaultPullRequestService) CreatePullRequest(pullRequest *dto.PullRequestDTO) (*dto.PullRequestDTO, error) {
	entity, err := mappers.PullRequestToEntity(pullRequest)
	if err != nil {
		return nil, err
	}

	pr, err := s.prDomainServ.CreateWithReviewers(entity)
	if err != nil {
		return nil, err
	}
	dto, err := mappers.PullRequestToDTO(pr)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *DefaultPullRequestService) MarkAsMerged(pullRequestID v.ID) (*dto.PullRequestDTO, error) {
	pr, err := s.prDomainServ.MarkAsMerged(pullRequestID)
	if err != nil {
		return nil, err
	}
	dto, err := mappers.PullRequestToDTO(pr)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *DefaultPullRequestService) ReassignReviewer(userID v.ID, pullRequestID v.ID) (*dto.PullRequestWithNewReviewerIDDTO, error) {
	newReviewer, err := s.prDomainServ.ReassignReviewer(userID, pullRequestID)
	if err != nil {
		return nil, err
	}
	d, err := mappers.PullRequestToDTO(&newReviewer.PullRequest)
	if err != nil {
		return nil, err
	}
	response := dto.PullRequestWithNewReviewerIDDTO{
		PullRequest:       d,
		NewReviewerUserID: newReviewer.NewReviewerID,
	}
	return &response, nil
}

func (s *DefaultPullRequestService) FindPullRequestsByReviewer(userID v.ID) ([]*dto.PullRequestDTO, error) {
	prs, err := s.prRepo.FindPullRequestsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	return mappers.PullRequestsToDTOs(prs)
}
