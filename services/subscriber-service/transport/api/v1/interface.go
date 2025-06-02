package v1

import (
	"context"

	"subscriber-service/pkg/id"
	svcmodel "subscriber-service/service/model"
	"subscriber-service/transport/api/v1/model"
)

type SubscriberService interface {
	SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error
	ConfirmSubscription(ctx context.Context, token string) (svcmodel.Subscription, error)
	UnsubscribeFromNewsletter(ctx context.Context, token string) error
	DeleteNewsletterSubscriptions(ctx context.Context, newsletterId id.Newsletter) error
	SendPublishedPost(ctx context.Context, post model.Post) error
}
