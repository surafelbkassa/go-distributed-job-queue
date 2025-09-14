package usecases

import (
	"errors"
	// "github.com/surafelbkassa/go-distributed-job-queue/Domain"
	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
	// "github.com/surafelbkassa/go-distributed-job-queue/Repository"
)

type UserUsecase struct {
	repo domain.IUserRepo
}

func NewUserUsecase(ur domain.IUserRepo) *UserUsecase {
	return &UserUsecase{
		repo: ur,
	}
}

func Register(uc *UserUsecase, user *domain.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing user data")
	}
	err := uc.repo.Register(user)
	if err != nil {
		return err
	}
	return nil
}
