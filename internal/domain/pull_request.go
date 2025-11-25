package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type PullRequestDTO struct {
	ID          ID
	Title       string
	AuthorID    ID
	CreatedAt   time.Time
	Status      string
	MergedAt    *time.Time
	ReviewerIDs []ID
}

const MaxReviewersCount int = 2

var (
	ErrPRAlreadyMerged           = errors.New("PR is already merged")
	ErrMaxReviewersCount         = fmt.Errorf("maximum number of reviewers is %d", MaxReviewersCount)
	ErrAlreadyAssignedAsReviewer = errors.New("user already assiggned as reviewer")
)

type PullRequest struct {
	id        ID
	title     string
	authorID  ID
	createdAt time.Time
	status    PRStatus
	mergedAt  *time.Time
	// slice (not map) because reviewers count is often not large
	reviewerIDs []ID
}

func NewPullRequest(title string, authorID ID) (*PullRequest, error) {
	return &PullRequest{
		NewID(),
		title,
		authorID,
		time.Now(),
		OPEN,
		nil,
		make([]ID, 0, MaxReviewersCount),
	}, nil
}

func NewExistingPullRequest(pullRequest *PullRequestDTO) (*PullRequest, error) {
	if pullRequest == nil {
		return nil, errors.New("dto cannot be nil")
	}

	status, err := NewPRStatusFromString(pullRequest.Status)
	if err != nil {
		return nil, err
	}

	if len(pullRequest.ReviewerIDs) > MaxReviewersCount {
		return nil, ErrMaxReviewersCount
	}

	err = validateIDsUniqueness(pullRequest.ReviewerIDs)
	if err != nil {
		return nil, fmt.Errorf("pull requests: %w", err)
	}

	if status == MERGED && pullRequest.MergedAt == nil {
		return nil, errors.New("PR marked as merged, but time is not specified")
	}
	if status == OPEN && pullRequest.MergedAt != nil {
		return nil, errors.New("PR marked as opened, but merege time not specified")
	}

	reviewerIDs := make([]ID, 0, MaxReviewersCount)
	reviewerIDs = append(reviewerIDs, pullRequest.ReviewerIDs...)

	return &PullRequest{
		pullRequest.ID,
		pullRequest.Title,
		pullRequest.AuthorID,
		pullRequest.CreatedAt,
		status,
		pullRequest.MergedAt,
		reviewerIDs,
	}, nil
}

func (p *PullRequest) ID() ID {
	return p.id
}

func (p *PullRequest) Title() string {
	return p.title
}

func (p *PullRequest) AuthorID() ID {
	return p.authorID
}

func (p *PullRequest) Status() PRStatus {
	return p.status
}

func (p *PullRequest) MergedAt() *time.Time {
	if p.mergedAt == nil {
		return nil
	}
	copy := *p.mergedAt
	return &copy
}

func (p *PullRequest) CreatedAt() time.Time {
	return p.createdAt
}

func (p *PullRequest) ReviewerIDs() []ID {
	return slices.Clone(p.reviewerIDs)
}

func (p *PullRequest) AssignReviewer(reviewerID ID) error {
	if len(p.reviewerIDs) == MaxReviewersCount {
		return ErrMaxReviewersCount
	}

	if p.status == MERGED {
		return ErrPRAlreadyMerged
	}

	if slices.Contains(p.reviewerIDs, reviewerID) {
		return fmt.Errorf("%w: id=%v", ErrAlreadyAssignedAsReviewer, reviewerID)
	}

	p.reviewerIDs = append(p.reviewerIDs, reviewerID)
	return nil
}

func (p *PullRequest) UnassignReviewer(reviewerID ID) error {
	if p.status == MERGED {
		return ErrPRAlreadyMerged
	}

	idx := slices.Index(p.reviewerIDs, reviewerID)
	if idx == -1 {
		return fmt.Errorf("no user with id=%d inside reviewers list", reviewerID)
	}

	p.reviewerIDs = slices.Delete(p.reviewerIDs, idx, idx+1)
	return nil
}

func (p *PullRequest) MarkAsMerged() {
	if p.status == MERGED {
		return
	}
	time := time.Now()
	p.status = MERGED
	p.mergedAt = &time
}

func validateIDsUniqueness(ids []ID) error {
	seen := make(map[ID]bool)

	for _, id := range ids {
		if _, exists := seen[id]; exists {
			return fmt.Errorf("at least one duplicated id=%v", id.String())
		}
		seen[id] = true
	}

	return nil
}
