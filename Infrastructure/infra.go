package infrastructure

import (
	"time"

	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Client, dbName string) *UserRepo {
	return &UserRepo{
		collection: db.Database(dbName).Collection("users"),
	}
}

func (r *UserRepo) Create(user *domain.User) error {
	user.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(nil, user)
	return err
}
