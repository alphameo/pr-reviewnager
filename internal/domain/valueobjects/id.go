// Package valueobjects provides value objects of domain model
package valueobjects

import (
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	id := ID(uuid.Must(uuid.NewV7()))
	return id
}

func NewIDFromString(str string) (ID, error) {
	value, err := uuid.Parse(str)
	if err != nil {
		return ID(uuid.Nil), err
	}
	id := ID(value)
	return id, nil
}

func (id ID) String() string {
	return uuid.UUID.String(uuid.UUID(id))
}
