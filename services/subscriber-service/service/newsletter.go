package service

import (
	"context"

	"subscriber-service/pkg/id"
)

func (s Service) DeleteNewsletterSubscriptions(ctx context.Context, newsletterId id.Newsletter) error {
	err := s.repository.DeleteNewsletterSubscriptions(ctx, newsletterId)
	if err != nil {
		return err
	}
	return err
}
