package rest

import (
	"context"

	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
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
	Authn(ctx context.Context, email, password string) (string, map[string]string, error)
	New(ctx context.Context, loggedInUserID, password string, user bizmodels.User) (string, error)
	Get(ctx context.Context, id string) (*bizmodels.User, error)
	GetAll(ctx context.Context) ([]*bizmodels.User, error)
	Edit(ctx context.Context, id, loggedInUserID string, updateCusKey bool, userEdit bizmodels.UserEdit) (*bizmodels.User, error)
	Delete(ctx context.Context, id, loggedInUserID string) error
}
