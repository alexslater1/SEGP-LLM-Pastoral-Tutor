package entity

const (
	UserEntityId    = "user"
	UserDescription = "This is the external user who gave the initial query. The task which you can give them should be a question which is refining their original query, if necessary."
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
