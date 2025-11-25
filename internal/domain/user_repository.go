package domain

type UserRepository interface {
	Repository[User, ID]
}
