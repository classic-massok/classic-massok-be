package resolvers

import (
	"context"

	"github.com/classic-massok/classic-massok-be/api/graphql/models"
)

func (m *mutation) Login(ctx context.Context, input models.LoginInput) (*models.User, error) {
	return nil, nil
}
