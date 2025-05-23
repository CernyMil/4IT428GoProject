package v1

import (
	"context"

	"subscriber-api/pkg/id"
	svcmodel "subscriber-api/service/model"
)

type SubscriberService interface {
	SubscribeToNewsletter(ctx context.Context, subReq svcmodel.SubscribeRequest) error
	ConfirmSubscription(ctx context.Context, subReq svcmodel.SubscribeRequest) (svcmodel.Subscription, error)
	UnsubscribeFromNewsletter(ctx context.Context, newsletterId id.Newsletter, subscriptionId id.Subscription) error
	DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error
	SendPublishedPost(ctx context.Context, post svcmodel.Post) error

	//SendConfirmationEmail(to, newsletterName, confirmLink string) error
	//SendNewsletterEmail(to, newsletterName, content, unsubscribeLink string) error
}
