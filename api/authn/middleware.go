package authn

import (
	"context"
	"time"

	"github.com/classic-massok/classic-massok-be/business"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserIDKey = "userID"
	RolesKey  = "roles"
)

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
			return validateToken(token, "api/authn/access_public_key.pem")
		})
		if err != nil {
			return a.validateRefreshToken(c, tokenString, next)
		}

		claims, ok := token.Claims.(*tokenClaims)
		if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != accessTokenType ||
			time.Now().After(time.Unix(claims.ExpiresAt, 0)) || claims.Issuer != "classic-massok.auth.service" {
			return next(c)
			// return fmt.Errorf("invalid token: authentication failed")
		}

		user, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID)
		if err != nil {
			return next(c)
			// return fmt.Errorf("invalid token: authentication failed")
		}

		c.Set(UserIDKey, user.GetID())
		c.Set(RolesKey, user.Roles)
		return next(c)
	}
}

func (a *AuthnMW) validateRefreshToken(c echo.Context, tokenString string, next echo.HandlerFunc) error {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, "api/authn/refresh_public_key.pem")
	})
	if err != nil {
		return next(c)
		// return err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != refreshTokenType || time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return next(c)
		// return fmt.Errorf("invalid token: authentication failed")
	}

	user, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID)
	if err != nil {
		return next(c)
		// return fmt.Errorf("invalid token: authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(claims.CusKey), []byte(user.GetCusKey())); err != nil {
		return next(c)
		// return fmt.Errorf("invalid token: authentication failed")
	}

	c.Set(UserIDKey, claims.UserID)
	c.Set(RolesKey, user.Roles)
	return next(c)
}

type userGetter interface {
	Get(ctx context.Context, id string) (*business.User, error)
}
