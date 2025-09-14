package repository

import domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"

type UserRepository interface {
	Register(user *domain.User) error
}

func Register(user *domain.User, repo UserRepository) error {
	if user == nil {
		return nil
	}
	return repo.Register(user)
}
