package entities

import (
	"errors"
	"fmt"
	"slices"
	"time"

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
	reviewerIDs []v.ID
}

func NewPullRequest(title string, authorID v.ID) *PullRequest {
	return NewExistingPullRequest(v.NewID(), title, authorID, time.Now(), v.OPEN, nil)
}

func NewExistingPullRequest(id v.ID, title string, authorID v.ID, createdAt time.Time, status v.PRStatus, mergedAt *time.Time) *PullRequest {
	p := PullRequest{
		id:          id,
		title:       title,
		authorID:    authorID,
		createdAt:   createdAt,
		status:      status,
		mergedAt:    mergedAt,
		reviewerIDs: make([]v.ID, 0, MaxCountOfReviewers),
	}

	return &p
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
