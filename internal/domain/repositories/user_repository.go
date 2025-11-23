package repositories

import (
	"github.com/alphameo/pr-reviewnager/internal/domain/dto"
	e "github.com/alphameo/pr-reviewnager/internal/domain/entities"
	v "github.com/alphameo/pr-reviewnager/internal/domain/valueobjects"
)

type UserRepository interface {
	Repository[e.User, dto.UserDTO, v.ID]
}
