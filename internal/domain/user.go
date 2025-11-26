package domain

type User struct {
	id     ID
	name   string
	active bool
}

func ExistingUser(
	id ID,
	name string,
	active bool,
) *User {
	return &User{
		id:     id,
		name:   name,
		active: active,
	}
}

func NewUser(name string, active bool) (*User, error) {
	return &User{
		id:     NewID(),
		name:   name,
		active: active,
	}, nil
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Active() bool {
	return u.active
}

func (u *User) SetActive(active bool) {
	u.active = active
}

func (u *User) Validate() error {
	// TODO: validation for name
	return nil
}
