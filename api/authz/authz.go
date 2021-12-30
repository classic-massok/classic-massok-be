package authz

import (
	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var ( // TODO: is this how we should declare errors? Should we have a centralized error repo?
	errNotFound     = errors.New("not found")
	errUnauthorized = errors.New("unauthorized")
	errForbidden    = errors.New("forbiddden")
	errServerError  = errors.New("server error")
)

func LoadResource(c echo.Context, bizRepo resourceGetter, resourceType, resourceID string) error {
	resource, err := bizRepo.Get(c.Request().Context(), resourceType, resourceID)
	if err != nil {
		// TODO: log this
		return errNotFound
	}

	c.Set("resource", resource)
	return nil
}

func RequiresPermission(c echo.Context, aclBiz accessAllower, action string) error {
	rolesVal := c.Get(authn.RolesKey)
	if rolesVal == nil {
		return errUnauthorized
	}

	roles := rolesVal.(business.Roles)
	resource := c.Get(resourceKey)

	if len(roles) == 0 {
		return errForbidden
	}

	allowed, err := aclBiz.AccessAllowed(c.Request().Context(), roles, resource, action)
	if err != nil {
		// TODO: log this
		return errors.Wrap(errServerError, err.Error()) // TODO: is this the right order of errors?
	}

	if !allowed {
		return errForbidden
	}

	return nil
}
