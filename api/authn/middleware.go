package authn

import (
	"context"
	"fmt"
	"time"

	"github.com/classic-massok/classic-massok-be/business"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const UserIDKey = "userID"

type AuthnMW struct {
	UsersBiz userGetter
}

func (a *AuthnMW) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := extractToken(c.Request())
		if tokenString == "" {
			return next(c)
		}

		token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return validateToken(token, "access_public_key.pem")
		})
		if err != nil {
			return a.validateRefreshToken(c, tokenString, next)
		}

		claims, ok := token.Claims.(*tokenClaims)
		if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != accessTokenType ||
			time.Now().After(time.Unix(claims.ExpiresAt, 0)) || claims.Issuer != "classic-massok.auth.service" {
			return fmt.Errorf("invalid token: authentication failed")
		}

		if _, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID); err != nil {
			return fmt.Errorf("invalid token: authentication failed")
		}

		c.Set(UserIDKey, claims.UserID)
		return next(c)
	}
}

func (a *AuthnMW) validateRefreshToken(c echo.Context, tokenString string, next echo.HandlerFunc) error {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, "refresh_public_key.pem")
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != refreshTokenType || time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return fmt.Errorf("invalid token: authentication failed")
	}

	user, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID)
	if err != nil {
		return fmt.Errorf("invalid token: authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(claims.CusKey), []byte(user.GetCusKey())); err != nil {
		return fmt.Errorf("invalid token: authentication failed")
	}

	c.Set(UserIDKey, claims.UserID)
	return next(c)
}

type userGetter interface {
	Get(ctx context.Context, id string) (*business.User, error)
}
