package authn

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

const (
	AccessTokenType  = "Access"
	RefreshTokenType = "Refresh"
)

func GenerateAccessToken(c echo.Context, userID string) (string, int64, error) {
	expiry := time.Now().Add(15 * time.Minute).Unix()
	claims := tokenClaims{
		UserID:    userID,
		TokenType: AccessTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
			Issuer:    "classic-massok.auth.service",
		},
	}

	token, err := generateToken(c.Get(lib.AccessTokenPrvKeyPathKey).(string), claims)
	return token, expiry, err
}

func GenerateRefreshToken(c echo.Context, userID, cusKey string) (string, int64, error) {
	hashedCusKey, err := bcrypt.GenerateFromPassword([]byte(cusKey), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, fmt.Errorf("error generating token: %w", err)
	}

	expiry := time.Now().Add(18 * time.Hour).Unix()
	claims := tokenClaims{
		UserID:    userID,
		TokenType: RefreshTokenType,
		CusKey:    string(hashedCusKey),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
			Issuer:    "classic-massok.auth.service",
		},
	}

	token, err := generateToken(c.Get(lib.RefreshTokenPrvKeyPathKey).(string), claims)
	return token, expiry, err
}

func generateToken(pemFilePath string, claims jwt.Claims) (string, error) {
	signBytes, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(signBytes, "local") // TODO: abstract password from this
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(signKey)
}

func validateToken(token *jwt.Token, pemFileName string) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("Unexpected signing method in auth token")
	}

	verifyBytes, err := ioutil.ReadFile(pemFileName)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return verifyKey, nil
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 || authHeaderContent[0] != "Bearer" {
		return ""
	}

	return authHeaderContent[1]
}

type tokenClaims struct {
	UserID    string `json:"userId"`
	CusKey    string `json:"cusKey,omitempty"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

//counterfeiter:generate . userGetter
type userGetter interface {
	Get(ctx context.Context, id string) (*bizmodels.User, error)
}

