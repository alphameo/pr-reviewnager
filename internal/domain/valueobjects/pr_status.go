package valueobjects

import (
	"fmt"
	"strings"
)

type PRStatus string

const (
	OPEN   PRStatus = "open"
	MERGED PRStatus = "merged"
)

func NewPRStatusFromString(value string) (PRStatus, error) {
	processed := strings.ToLower(strings.TrimSpace(value))
	switch processed {
	case "", "open", "opened":
		return OPEN, nil
	case "merged":
		return MERGED, nil
	default:
		return PRStatus(""), fmt.Errorf("no PR status: %s", processed)
	}
}

func (s PRStatus) String() string {
	return string(s)
}
