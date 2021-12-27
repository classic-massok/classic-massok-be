package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) Login(ctx context.Context, input models.LoginInput) (*models.AuthOutput, error) {
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

	return &models.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}

func (m *mutation) RefreshToken(ctx context.Context) (*models.AuthOutput, error) {
	c, err := echoContextFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	userID := c.Get(authn.UserIDKey).(string)

	bizUser, err := m.UsersBiz.Edit(ctx, userID, userID, true, business.UserEdit{})
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	accessToken, accessTokenExpiry, err := authn.GenerateAccessToken(bizUser.GetID())
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := authn.GenerateRefreshToken(userID, bizUser.GetCusKey())
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	return &models.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}
