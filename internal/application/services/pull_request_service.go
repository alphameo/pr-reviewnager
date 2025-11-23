package services

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/application/mappers"
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	r "github.com/alphameo/pr-reviewnager/internal/domain/repositories"
	ds "github.com/alphameo/pr-reviewnager/internal/domain/services"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type PullRequestService interface {
	CreatePullRequest(pullRequest *dto.PullRequestDTO) (*dto.PullRequestDTO, error)
	MarkAsMerged(pullRequestID v.ID) (*dto.PullRequestDTO, error)
	ReassignReviewer(userID v.ID, pullRequestID v.ID) (*PullRequestWithNewReviewerIDDTO, error)
	FindPullRequestsByReviewer(userID v.ID) ([]*dto.PullRequestDTO, error)
}

type PullRequestWithNewReviewerIDDTO struct {
	PullRequest       *dto.PullRequestDTO
	NewReviewerUserID v.ID
}

type DefaultPullRequestService struct {
	prDomainServ ds.PullRequestDomainService
	prRepo       r.PullRequestRepository
}

func NewDefaultPullRequestService(
	pullRequestDomainService ds.PullRequestDomainService,
	pullRequestRepository r.PullRequestRepository,
) (*DefaultPullRequestService, error) {
	if pullRequestDomainService == nil {
		return nil, errors.New("pullRequestDomainService cannot bi nil")
	}
	if pullRequestRepository == nil {
		return nil, errors.New("PullRequestRepository cannot be nil")
	}

	return &DefaultPullRequestService{
		prDomainServ: pullRequestDomainService,
		prRepo:       pullRequestRepository,
	}, nil
}

func (s *DefaultPullRequestService) CreatePullRequest(pullRequest *dto.PullRequestDTO) (*dto.PullRequestDTO, error) {
	entity, err := mappers.PullRequestToEntity(pullRequest)
	if err != nil {
		return nil, err
	}

	pr, err := s.prDomainServ.CreateWithReviewers(entity)
	if errors.Is(err, ds.ErrAuthorNotFound) || errors.Is(err, ds.ErrTeamNotFound) {
	} else if errors.Is(err, ds.ErrPRAlreadyExists) {
		return nil, ErrPRExists
	} else if err != nil {
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
	if errors.Is(err, ds.ErrPRNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	dto, err := mappers.PullRequestToDTO(pr)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *DefaultPullRequestService) ReassignReviewer(userID v.ID, pullRequestID v.ID) (*PullRequestWithNewReviewerIDDTO, error) {
	newReviewer, err := s.prDomainServ.ReassignReviewer(userID, pullRequestID)
	if errors.Is(err, ds.ErrPRNotFound) || errors.Is(err, ds.ErrUserNotFound) {
		return nil, ErrNotFound
	} else if errors.Is(err, ds.ErrPRAlreadyMerged) {
		return nil, ErrPRAlreadyMerged
	} else if errors.Is(err, ds.ErrUserNotReviewer) {
		return nil, ErrNotAssigned
	} else if errors.Is(err, ds.ErrNoReviewCandidates) {
		return nil, ErrNoCandidate
	} else if err != nil {
		return nil, err
	}
	d, err := mappers.PullRequestToDTO(&newReviewer.PullRequest)
	if err != nil {
		return nil, err
	}

	return &PullRequestWithNewReviewerIDDTO{
		PullRequest:       d,
		NewReviewerUserID: newReviewer.NewReviewerID,
	}, nil
}

func (s *DefaultPullRequestService) FindPullRequestsByReviewer(userID v.ID) ([]*dto.PullRequestDTO, error) {
	prs, err := s.prRepo.FindPullRequestsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
