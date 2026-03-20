package infrastructure

import (
	"sync"

	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
)

type InMemoryUserRepository struct {
	users []*domain.User
	mu    sync.Mutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make([]*domain.User, 0),
	}
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, user)
	return nil
}

// Optionally, add methods like FindByEmail, etc.
