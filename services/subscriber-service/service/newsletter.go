package service

import (
	"context"

	"subscriber-service/pkg/id"
)

/*
func (s Service) CreateNewsletter(ctx context.Context, newsletterId id.Newsletter) error {
	err := s.repository.CreateNewsletter(ctx, newsletterId)
	if err != nil {
		return err
	}
	return err
}
*/

func (s Service) DeleteNewsletterSubscriptions(ctx context.Context, newsletterId id.Newsletter) error {
	err := s.repository.DeleteNewsletterSubscriptions(ctx, newsletterId)
	if err != nil {
		return err
	}
	return err
}
