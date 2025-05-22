package service

import (
	"context"
	"subscriber-api/pkg/id"
	"subscriber-api/service/model"
	svcmodel "subscriber-api/service/model"
)

type Repository interface {
	AddSubscription(ctx context.Context, subscription svcmodel.Subscription) (*model.Subscription, error)
	RemoveSubscription(ctx context.Context, newsletter id.Newsletter, email string) error
	//ReadUser(ctx context.Context, userID id.User) (*model.User, error)
	//ListUser(ctx context.Context) ([]model.User, error)
}

type Service struct {
	repository Repository
}

func NewService(
	repository Repository,
) (Service, error) {
	return Service{
		repository: repository,
	}, nil
}
