package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

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
	return NewPullRequestWithID(NewID(), title, authorID)
}

func NewPullRequestWithID(id ID, title string, authorID ID) (*PullRequest, error) {
	return &PullRequest{
		id,
		title,
		authorID,
		time.Now(),
		PROpen,
		nil,
		make([]ID, 0, MaxReviewersCount),
	}, nil
}

func ExistingPullRequest(
	id ID,
	title string,
	authorID ID,
	createdAt time.Time,
	status PRStatus,
	mergedAt *time.Time,
	reviewerIDs []ID,
) *PullRequest {
	rIDs := make([]ID, 0, MaxReviewersCount)
	rIDs = append(rIDs, reviewerIDs...)

	return &PullRequest{
		id,
		title,
		authorID,
		createdAt,
		status,
		mergedAt,
		rIDs,
	}
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

	if p.status == PRMerged {
		return ErrPRAlreadyMerged
	}

	if slices.Contains(p.reviewerIDs, reviewerID) {
		return fmt.Errorf("%w: id=%v", ErrAlreadyAssignedAsReviewer, reviewerID)
	}

	p.reviewerIDs = append(p.reviewerIDs, reviewerID)
	return nil
}

func (p *PullRequest) UnassignReviewer(reviewerID ID) error {
	if p.status == PRMerged {
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
	if p.status == PRMerged {
		return
	}
	time := time.Now()
	p.status = PRMerged
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

func (p *PullRequest) Validate() error {
	if len(p.reviewerIDs) > MaxReviewersCount {
		return ErrMaxReviewersCount
	}

	err := validateIDsUniqueness(p.reviewerIDs)
	if err != nil {
		return fmt.Errorf("pull requests: %w", err)
	}

	if p.status == PRMerged && p.mergedAt == nil {
		return errors.New("PR marked as merged, but time is not specified")
	}
	if p.status == PROpen && p.mergedAt != nil {
		return errors.New("PR marked as opened, but merege time not specified")
	}

	return nil
}
