package service

import (
	"context"
	"subscriber-service/pkg/id"
	dbmodel "subscriber-service/repository/model"
	svcmodel "subscriber-service/service/model"
)

type Repository interface {
	AddSubscription(ctx context.Context, subscription svcmodel.Subscription) error
	DeleteSubscription(ctx context.Context, unsubReq svcmodel.UnsubscribeRequest) error
	GetSubscribers(ctx context.Context, newsletterId id.Newsletter) ([]dbmodel.SubscriberInfo, error)
	DeleteNewsletterSubscriptions(ctx context.Context, newsletterId id.Newsletter) error
	GetNewsletterById(ctx context.Context, newsletterId id.Newsletter) (id.Newsletter, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) (Service, error) {
	return Service{
		repository: repository,
	}, nil
}
