package usecase

import (
	"context"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/domain"
)

type HandlerFunc func(ctx context.Context, j *domain.Job) error

type Registry struct{ m map[string]HandlerFunc }

func NewRegistry() *Registry { return &Registry{m: map[string]HandlerFunc{}} }

func (r *Registry) Register(jobType string, h HandlerFunc) { r.m[jobType] = h }

func (r *Registry) Lookup(jobType string) (HandlerFunc, bool) {
	h, ok := r.m[jobType]
	return h, ok
}
