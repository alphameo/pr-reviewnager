// Package repositories provides interfaces for real repositories
package repositories

type Repository[T any, ID any] interface {
	Create(entity *T) error
	FindByID(id ID) (*T, error)
	FindAll() ([]T, error)
	Update(entity *T) error
	DeleteByID(id ID) error
}
