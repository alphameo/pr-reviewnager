package app

import (
	"errors"

	"github.com/alphameo/pr-reviewnager/internal/domain"
)

type PullRequestService interface {
	CreatePullRequest(pullRequest *domain.PullRequestDTO) (*domain.PullRequestDTO, error)
	MarkAsMerged(pullRequestID domain.ID) (*domain.PullRequestDTO, error)
	ReassignReviewer(userID domain.ID, pullRequestID domain.ID) (*PullRequestWithNewReviewerIDDTO, error)
	FindPullRequestsByReviewer(userID domain.ID) ([]*domain.PullRequestDTO, error)
}

type PullRequestWithNewReviewerIDDTO struct {
	PullRequest       *domain.PullRequestDTO
	NewReviewerUserID domain.ID
}

type DefaultPullRequestService struct {
	prDomainServ domain.PullRequestDomainService
	prRepo       domain.PullRequestRepository
}

func NewDefaultPullRequestService(
	pullRequestDomainService domain.PullRequestDomainService,
	pullRequestRepository domain.PullRequestRepository,
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

func (s *DefaultPullRequestService) CreatePullRequest(pullRequest *domain.PullRequestDTO) (*domain.PullRequestDTO, error) {
	entity, err := PullRequestToEntity(pullRequest)
	if err != nil {
		return nil, err
	}

	pr, err := s.prDomainServ.CreateWithReviewers(entity)
	if errors.Is(err, domain.ErrAuthorNotFound) || errors.Is(err, domain.ErrTeamNotFound) {
	} else if errors.Is(err, domain.ErrPRAlreadyExists) {
		return nil, ErrPRExists
	} else if err != nil {
		return nil, err
	}
	dto, err := PullRequestToDTO(pr)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *DefaultPullRequestService) MarkAsMerged(pullRequestID domain.ID) (*domain.PullRequestDTO, error) {
	pr, err := s.prDomainServ.MarkAsMerged(pullRequestID)
	if errors.Is(err, domain.ErrPRNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	dto, err := PullRequestToDTO(pr)
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *DefaultPullRequestService) ReassignReviewer(userID domain.ID, pullRequestID domain.ID) (*PullRequestWithNewReviewerIDDTO, error) {
	newReviewer, err := s.prDomainServ.ReassignReviewer(userID, pullRequestID)
	if errors.Is(err, domain.ErrPRNotFound) || errors.Is(err, domain.ErrUserNotFound) {
		return nil, ErrNotFound
	} else if errors.Is(err, domain.ErrPRAlreadyMerged) {
		return nil, ErrPRAlreadyMerged
	} else if errors.Is(err, domain.ErrUserNotReviewer) {
		return nil, ErrNotAssigned
	} else if errors.Is(err, domain.ErrNoReviewCandidates) {
		return nil, ErrNoCandidate
	} else if err != nil {
		return nil, err
	}
	d, err := PullRequestToDTO(&newReviewer.PullRequest)
	if err != nil {
		return nil, err
	}

	return &PullRequestWithNewReviewerIDDTO{
		PullRequest:       d,
		NewReviewerUserID: newReviewer.NewReviewerID,
	}, nil
}

func (s *DefaultPullRequestService) FindPullRequestsByReviewer(userID domain.ID) ([]*domain.PullRequestDTO, error) {
	prs, err := s.prRepo.FindPullRequestsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
