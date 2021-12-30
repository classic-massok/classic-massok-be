package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	graphqlmodels "github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) Login(ctx context.Context, input graphqlmodels.LoginInput) (*graphqlmodels.AuthOutput, error) {
	userID, cusKey, err := m.UsersBiz.Authn(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	_, err = m.UsersBiz.Edit(ctx, userID, userID, true, business.UserEdit{})
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	accessToken, accessTokenExpiry, err := authn.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := authn.GenerateRefreshToken(userID, cusKey)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	c, err := echoContextFromContext(ctx)
	if err != nil {
		// log error here
	} else {
		c.Set(authn.UserIDKey, userID)
	}

	return &graphqlmodels.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}

func (m *mutation) RefreshToken(ctx context.Context) (*graphqlmodels.AuthOutput, error) {
	c, err := echoContextFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	tokenType := c.Get(authn.TokenTypeKey).(string)
	if tokenType != authn.RefreshTokenType {
		return nil, fmt.Errorf("forbidden")
	}

	userID := c.Get(authn.UserIDKey).(string) // TODO: is this safe?

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

	return &graphqlmodels.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}
