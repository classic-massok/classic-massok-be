package authz

import (
	"context"
	"fmt"

	bizmodels "github.com/classic-massok/classic-massok-be/business/models"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

func LoadResource(c echo.Context, bizRepo resourceGetter, resourceType, resourceID string) error {
	resource, err := bizRepo.Get(c.Request().Context(), resourceType, resourceID)
	if err != nil {
		// TODO: log this
		return lib.ErrNotFound
	}

	c.Set(resourceKey, resource)
	return nil
}

func RequiresPermission(c echo.Context, aclBiz accessAllower, action string) error {
	rolesVal := c.Get(lib.RolesKey)
	if rolesVal == nil {
		return lib.ErrUnauthorized
	}
	roles := rolesVal.(bizmodels.Roles)

	userIDVal := c.Get(lib.UserIDKey)
	if userIDVal == nil {
		return lib.ErrUnauthorized
	}
	userID := userIDVal.(string)

	resource := c.Get(resourceKey)

	if len(roles) == 0 {
		return lib.ErrForbidden
	}

	allowed, err := aclBiz.AccessAllowed(c.Request().Context(), resource, action, userID, roles)
	if err != nil {
		// TODO: log this
		return fmt.Errorf("%w: %v", lib.ErrServerError, err) // TODO: is this the right order of errors?
	}

	if !allowed {
		return lib.ErrForbidden
	}

	return nil
}

//counterfeiter:generate . resourceGetter
type resourceGetter interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

//counterfeiter:generate . accessAllower
type accessAllower interface {
	AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles bizmodels.Roles) (bool, error)
}
