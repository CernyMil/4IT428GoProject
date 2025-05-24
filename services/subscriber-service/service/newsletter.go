package service

import (
	"context"

	"subscriber-api/pkg/id"
)

func (s Service) DeleteNewsletter(ctx context.Context, newsletterId id.Newsletter) error {
	err := s.repository.DeleteNewsletter(ctx, newsletterId)
	if err != nil {
		return err
	}
	return err
}
