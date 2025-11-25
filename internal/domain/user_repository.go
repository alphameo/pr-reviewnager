package domain

type UserRepository interface {
	Repository[User, UserDTO, ID]
}
