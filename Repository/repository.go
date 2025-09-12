package repository

import domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"

type UserRepository interface {
	Save(user *domain.User) error
}