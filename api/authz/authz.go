package authz

import (
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

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
	roles := rolesVal.(business.Roles)

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
		return errors.Wrap(lib.ErrServerError, err.Error()) // TODO: is this the right order of errors?
	}

	if !allowed {
		return lib.ErrForbidden
	}

	return nil
}
