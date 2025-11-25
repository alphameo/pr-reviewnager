package domain

import (
	"fmt"
	"strings"
)

type PRStatus string

const (
	PROpen   PRStatus = "open"
	PRMerged PRStatus = "merged"
)

func NewPRStatus(value string) (PRStatus, error) {
	processed := strings.ToLower(strings.TrimSpace(value))
	switch processed {
	case "", "open", "opened":
		return PROpen, nil
	case "merged":
		return PRMerged, nil
	default:
		return PRStatus(""), fmt.Errorf("no PR status: %s", processed)
	}
}

func ExistingPRStatus(value string) PRStatus {
	return PRStatus(value)
}

func (s PRStatus) String() string {
	return string(s)
}
