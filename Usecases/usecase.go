package usecases

import (
	"errors"

	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
	repository "github.com/surafelbkassa/go-distributed-job-queue/Repository"
)

type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

type RegisterUserUsecase struct {
	UserRepo repository.UserRepository
}

func NewRegisterUserUsecase(userRepo repository.UserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{UserRepo: userRepo}
}

func (uc *RegisterUserUsecase) Execute(input RegisterUserInput) (*domain.User, error) {
	// Add validation as needed
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("all fields are required")
	}

	user := domain.NewUser(input.Username, input.Email, input.Password)
	err := uc.UserRepo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
