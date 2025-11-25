// Package dto provides data transfer objects as animic domain entities
package dto

import v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"

type User struct {
	ID     v.ID
	Name   string
	Active bool
}
