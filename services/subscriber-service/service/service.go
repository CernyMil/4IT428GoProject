package service

import (
	"context"
	"subscriber-api/pkg/id"
	dbmodel "subscriber-api/repository/model"
	svcmodel "subscriber-api/service/model"
)

type Repository interface {
	AddSubscription(ctx context.Context, newsletterId id.Newsletter, subscribtionId id.Subscription, email string, token string) (*svcmodel.Subscription, error)
	DeleteSubscription(ctx context.Context, newsletterId id.Newsletter, subscribtionId id.Subscription) error
	GetSubscribers(ctx context.Context, newsletterId id.Newsletter) ([]dbmodel.SubscriberInfo, error)
	DeleteNewsletter(ctx context.Context, newsletter id.Newsletter) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) (Service, error) {
	return Service{
		repository: repository,
	}, nil
}
