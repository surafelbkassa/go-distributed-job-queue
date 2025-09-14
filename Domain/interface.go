package domain

import "context"

type IUserRepo interface {
	Register(user *User) error
}

type IUserUsecase interface {
	RegisterUser(ctx context.Context, user *User) (*User, error)
}
