// Package dto provides data transfer objects for application layer
package dto

import v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"

type UserDTO struct {
	ID     v.ID
	Name   string
	Active bool
}

type UserWithTeamNameDTO struct {
	User   *UserDTO
	TeamName string
}
