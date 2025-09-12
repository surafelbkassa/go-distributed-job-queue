package usecases

import (
    "errors"
    "github.com/surafelbkassa/go-distributed-job-queue/Domain"
    "github.com/surafelbkassa/go-distributed-job-queue/Repository"
)

type RegisterUserInput struct {
    Username string
    Email    string
    Password string
}

type RegisterUserUsecase struct {
    UserRepo Repository.UserRepository
}

func NewRegisterUserUsecase(userRepo Repository.UserRepository) *RegisterUserUsecase {
    return &RegisterUserUsecase{UserRepo: userRepo}
}

func (uc *RegisterUserUsecase) Execute(input RegisterUserInput) (*Domain.User, error) {
    // Add validation as needed
    if input.Username == "" || input.Email == "" || input.Password == "" {
        return nil, errors.New("all fields are required")
    }

    user := Domain.NewUser(input.Username, input.Email, input.Password)
    err := uc.UserRepo.Save(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}