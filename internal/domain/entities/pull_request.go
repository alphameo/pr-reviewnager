package entities

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

const MaxReviewersCount int = 2

var (
	ErrPRAlreadyMerged           = errors.New("PR is already MERGED")
	ErrMaxReviewersCount         = fmt.Errorf("maximum number of reviewers is %d", MaxReviewersCount)
	ErrAlreadyAssignedAsReviewer = errors.New("user already assiggned as reviewer")
)

type PullRequest struct {
	id          v.ID
	title       string
	authorID    v.ID
	createdAt   time.Time
	status      v.PRStatus
	mergedAt    *time.Time
	// slice (not map) because reviewers count is often not large
	reviewerIDs []v.ID
}

func NewPullRequest(title string, authorID v.ID) (*PullRequest, error) {
	return &PullRequest{
		v.NewID(),
		title,
		authorID,
		time.Now(),
		v.OPEN,
		nil,
		make([]v.ID, 0, MaxReviewersCount),
	}, nil
}

func NewExistingPullRequest(pullRequest *dto.PullRequestDTO) (*PullRequest, error) {
	if pullRequest == nil {
		return nil, errors.New("dto cannot be nil")
	}

	status, err := v.NewPRStatusFromString(pullRequest.Status)
	if err != nil {
		return nil, err
	}

	if len(pullRequest.ReviewerIDs) > MaxReviewersCount {
		return  nil, ErrMaxReviewersCount
	}

	err = validateIDsUniqueness(pullRequest.ReviewerIDs)
	if err != nil {
		return nil, fmt.Errorf("pull requests: %w", err)
	}

	if status == v.MERGED && pullRequest.MergedAt == nil {
		return nil, errors.New("PR marked as merged, but time is not specified")
	}
	if status == v.OPEN && pullRequest.MergedAt != nil {
		return nil, errors.New("PR marked as opened, but merege time not specified")
	}

	reviewerIDs := make([]v.ID, 0, MaxReviewersCount)
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

func (p *PullRequest) ID() v.ID {
	return p.id
}

func (p *PullRequest) Title() string {
	return p.title
}

func (p *PullRequest) AuthorID() v.ID {
	return p.authorID
}

func (p *PullRequest) Status() v.PRStatus {
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

func (p *PullRequest) ReviewerIDs() []v.ID {
	return slices.Clone(p.reviewerIDs)
}

func (p *PullRequest) NeedMoreReviewers() bool {
	return len(p.reviewerIDs) < MaxCountOfReviewers
}

func (p *PullRequest) AssignReviewer(reviewerID v.ID) error {
	if len(p.reviewerIDs) == MaxReviewersCount {
		return ErrMaxReviewersCount
	}

	if p.status == v.MERGED {
		return ErrPRAlreadyMerged
	}

	if slices.Contains(p.reviewerIDs, reviewerID) {
		return fmt.Errorf("%w: id=%v", ErrAlreadyAssignedAsReviewer, reviewerID)
	}

	p.reviewerIDs = append(p.reviewerIDs, reviewerID)
	return nil
}

func (p *PullRequest) UnassignReviewer(reviewerID v.ID) error {
	if p.status == v.MERGED {
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
	if p.status == v.MERGED {
		return
	}
	time := time.Now()
	p.status = v.MERGED
	p.mergedAt = &time
}

func validateIDsUniqueness(ids []v.ID) error {

	seen := make(map[v.ID]bool)

	for _, id := range ids {
		if _, exists := seen[id]; exists {
			return fmt.Errorf("at least one duplicated id=%v", id.String())
		}
		seen[id] = true
	}

	return nil
}
