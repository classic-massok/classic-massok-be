package resolvers

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	graphqlmodels "github.com/classic-massok/classic-massok-be/api/graphql/models"
	"github.com/classic-massok/classic-massok-be/business"
)

func (m *mutation) Login(ctx context.Context, input graphqlmodels.LoginInput) (*graphqlmodels.AuthOutput, error) {
	c, err := echoContextFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}

	userID, cusKeys, err := m.UsersBiz.Authn(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	c.Set(authn.UserIDKey, userID) // TODO: Do we need to do this? maybe for loggiing?
	ctx = context.WithValue(ctx, authn.CusKeysKey, cusKeys)

	bizUser, err := m.UsersBiz.Edit(ctx, userID, userID, true, business.UserEdit{})
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}

	accessToken, accessTokenExpiry, err := authn.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := authn.GenerateRefreshToken(userID, bizUser.GetCusKey(ctx.Value("IPAddress").(string)))
	if err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}

	return &graphqlmodels.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}

func (m *mutation) RefreshToken(ctx context.Context) (*graphqlmodels.AuthOutput, error) {
	c, err := echoContextFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	tokenType := c.Get(authn.TokenTypeKey).(string)
	if tokenType != authn.RefreshTokenType {
		return nil, fmt.Errorf("forbidden")
	}

	userID := c.Get(authn.UserIDKey).(string) // TODO: is this safe?
	ctx = context.WithValue(ctx, authn.CusKeysKey, c.Get(authn.CusKeysKey))

	bizUser, err := m.UsersBiz.Edit(ctx, userID, userID, true, business.UserEdit{})
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	accessToken, accessTokenExpiry, err := authn.GenerateAccessToken(bizUser.GetID())
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	refreshToken, refreshTokenExpiry, err := authn.GenerateRefreshToken(userID, bizUser.GetCusKey(c.Echo().IPExtractor(c.Request())))
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	return &graphqlmodels.AuthOutput{
		accessToken, accessTokenExpiry,
		refreshToken, refreshTokenExpiry,
	}, nil
}
