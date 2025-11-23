// Package repositories provides interfaces for real repositories
package repositories

type Repository[T any, DTO any, ID any] interface {
	Create(entity *T) error
	FindByID(id ID) (*DTO, error)
	FindAll() ([]*DTO, error)
	Update(entity *T) error
	DeleteByID(id ID) error
}
