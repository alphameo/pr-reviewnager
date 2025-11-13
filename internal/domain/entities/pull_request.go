package entities

import (
	"errors"
	"fmt"
	"slices"

	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

const MaxCountOfReviewers int = 2

type PullRequest struct {
	id          v.ID
	title       string
	authorID    v.ID
	status      v.PRStatus
	reviewerIDs []v.ID
}

func NewPullRequest(title string, authorID v.ID) *PullRequest {
	p := PullRequest{
		id:          v.NewID(),
		title:       title,
		authorID:    authorID,
		status:      v.OPEN,
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
	p.status = v.MERGED
}
