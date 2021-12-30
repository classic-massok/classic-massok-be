package rest

import (
	"context"

	"github.com/classic-massok/classic-massok-be/business"
	"github.com/labstack/echo/v4"
)

type Rest struct {
	AuthzMW  authzMW
	UsersBiz usersBiz
}

func (r *Rest) Configure(rest *echo.Group) {
	rest.GET("/hello-world", helloWorld)
}

func helloWorld(c echo.Context) error {
	return c.JSON(200, "Hello, World!")
}

type authzMW interface {
	LoadResource(resourceType string, resourceID string) echo.MiddlewareFunc
	RequiresPermission(action string) echo.MiddlewareFunc
}

type usersBiz interface {
	Authn(ctx context.Context, email, password string) (string, string, error)
	New(ctx context.Context, loggedInUserID, password string, user business.User) (string, error)
	Get(ctx context.Context, id string) (*business.User, error)
	GetAll(ctx context.Context) ([]*business.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit business.UserEdit) (*business.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
