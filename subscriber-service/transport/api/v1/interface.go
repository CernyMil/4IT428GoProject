package v1

import (
	"context"

	"subscriber-api/pkg/id"
	svcmodel "subscriber-api/service/model"
	"subscriber-api/transport/api/v1/model"
)

type SubscriberService interface {
	SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error
	ConfirmSubscription(ctx context.Context, subReq svcmodel.SubscribeRequest) (svcmodel.Subscription, error)
	UnsubscribeFromNewsletter(ctx context.Context, unsubReq svcmodel.UnsubscribeRequest) error
	DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error
	SendPublishedPost(ctx context.Context, post model.Post) error
}
