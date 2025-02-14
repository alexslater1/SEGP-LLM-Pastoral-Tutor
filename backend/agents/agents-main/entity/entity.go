package entity

const (
	UserEntityId    = "user"
	UserDescription = "The external user who gave the query"
)

type Entity interface {
	Id() string
	Description() string
}

type UserEntity struct {
}

func (u *UserEntity) Id() string {
	return UserEntityId
}

func (u *UserEntity) Description() string {
	return UserDescription
}
