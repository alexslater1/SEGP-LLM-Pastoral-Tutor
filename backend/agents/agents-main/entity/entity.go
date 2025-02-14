package entity

type Entity interface {
	Id() string
	Description() string
}

type UserEntity struct {
}

func (u *UserEntity) Id() string {
	return "user"
}

func (u *UserEntity) Description() string {
	return "The external user who gave the query"
}
