package v1

import (
	"context"

	"subscriber-service/pkg/id"
	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/api/v1/model"
)

type SubscriberService interface {
	SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error
	ConfirmSubscription(ctx context.Context, subReq svcmodel.SubscribeRequest) (svcmodel.Subscription, error)
	UnsubscribeFromNewsletter(ctx context.Context, unsubReq svcmodel.UnsubscribeRequest) error
	DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error
	CreateNewsletter(ctx context.Context, newsletterId id.Newsletter) error
	SendPublishedPost(ctx context.Context, post model.Post) error
}
