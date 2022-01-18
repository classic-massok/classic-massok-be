package authn

import (
	"time"

	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var parser = jwt.Parser{
	SkipClaimsValidation: true, // we validate ourselves for ease of flow
}

type AuthnMW struct {
	UsersBiz userGetter
}

func (a *AuthnMW) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := extractToken(c.Request())
		if tokenString == "" {
			return next(c)
		}

		token, err := parser.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return validateToken(token, c.Get(lib.AccessTokenPubKeyPathKey).(string))
		})
		if err != nil {
			return a.validateRefreshToken(c, tokenString, next)
		}

		claims, ok := token.Claims.(*tokenClaims)
		if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != AccessTokenType ||
			time.Now().After(time.Unix(claims.ExpiresAt, 0)) || claims.Issuer != "classic-massok.auth.service" {
			return next(c)
		}

		user, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID)
		if err != nil {
			return next(c)
		}

		c.Set(lib.UserIDKey, user.ID)
		c.Set(lib.RolesKey, user.Roles)
		c.Set(lib.CusKeysKey, user.CusKeys)
		c.Set(lib.TokenTypeKey, claims.TokenType)
		return next(c)
	}
}

func (a *AuthnMW) validateRefreshToken(c echo.Context, tokenString string, next echo.HandlerFunc) error {
	token, err := parser.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, c.Get(lib.RefreshTokenPubKeyPathKey).(string))
	})
	if err != nil {
		return next(c)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.TokenType != RefreshTokenType || time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return next(c)
	}

	user, err := a.UsersBiz.Get(c.Request().Context(), claims.UserID)
	if err != nil {
		return next(c)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(claims.CusKey), []byte(user.CusKeys[c.Echo().IPExtractor(c.Request())])); err != nil {
		return next(c)
	}

	c.Set(lib.UserIDKey, user.ID)
	c.Set(lib.TokenTypeKey, claims.TokenType)
	c.Set(lib.CusKeysKey, user.CusKeys)
	return next(c)
}

