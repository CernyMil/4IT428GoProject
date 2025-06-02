package v1

import (
	"context"
	id "newsletter-service/pkg/id"
	svcmodel "newsletter-service/service/model"
)

type NewsletterService interface {
	CreateNewsletter(ctx context.Context, input svcmodel.CreateNewsletterInput) (*svcmodel.Newsletter, error)
	ListNewsletters(ctx context.Context) ([]svcmodel.Newsletter, error)
	UpdateNewsletter(ctx context.Context, id id.Newsletter, input svcmodel.UpdateNewsletterInput) (*svcmodel.Newsletter, error)
	DeleteNewsletter(ctx context.Context, id id.Newsletter) error

	CreatePost(ctx context.Context, newsletterID id.Newsletter, input svcmodel.CreatePostInput) (*svcmodel.Post, error)
	ListPosts(ctx context.Context, newsletterID id.Newsletter) ([]svcmodel.Post, error)
	UpdatePost(ctx context.Context, newsletterID id.Newsletter, postID id.Post, input svcmodel.UpdatePostInput) (*svcmodel.Post, error)
	DeletePost(ctx context.Context, newsletterID id.Newsletter, postID id.Post) error
	PublishPost(ctx context.Context, newsletterID id.Newsletter, postID id.Post) (*svcmodel.Post, error)
}
