package entities

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

const MaxCountOfReviewers int = 2

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
	if !p.NeedMoreReviewers() {
		return fmt.Errorf("maximum number of reviewers is %d", MaxCountOfReviewers)
	}
	if p.status == v.MERGED {
		return errors.New("cannot assign new reviewer: PR is already MERGED")
	}

	if slices.Contains(p.reviewerIDs, reviewerID) {
		return fmt.Errorf("user with id=%v already assigned as review", reviewerID)
	}

	p.reviewerIDs = append(p.reviewerIDs, reviewerID)
	return nil
}

func (p *PullRequest) UnassignReviewer(reviewerID v.ID) error {
	if p.status == v.MERGED {
		return errors.New("cannot unassign reviewer: PR is already MERGED")
	}
	idx := slices.Index(p.reviewerIDs, reviewerID)
	if idx == -1 {
		return fmt.Errorf("cannot unassing reviwer: no user with id=%d inside reviewers list", reviewerID)
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
