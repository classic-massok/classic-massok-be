package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/api/graphql/models"
)

func (m *mutation) Login(ctx context.Context, input models.LoginInput) (*models.LoginOutput, error) {
	userID, cusKey, err := m.UsersBiz.Authn(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenExpiry, err := authn.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := authn.GenerateRefreshToken(userID, cusKey)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	return &models.LoginOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}
