package entity

const (
	UserEntityId    = "user"
	UserDescription = "The external user who gave the query"
)

var UserEntity = &userEntity{}

type Entity interface {
	Id() string
	Description() string
}

type userEntity struct {
}

func (u *userEntity) Id() string {
	return UserEntityId
}

func (u *userEntity) Description() string {
	return UserDescription
}
