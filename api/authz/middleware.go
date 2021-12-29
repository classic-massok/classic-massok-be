package authz

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/authn"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/labstack/echo/v4"
)

const resourceKey = "resource"

type AuthzMW struct {
	ACLBiz  accessAllower
	BizRepo resourceGetter
}

func (a *AuthzMW) LoadResource(resourceType string, resourceID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resource, err := a.BizRepo.Get(c.Request().Context(), resourceType, resourceID)
			if err != nil {
				return err // TODO: do actual json error response here (404)
			}

			c.Set("resource", resource)
			return next(c)
		}
	}
}

func (a *AuthzMW) RequiresPermission(action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roles := c.Get(authn.RolesKey).(business.Roles)
			resource := c.Get(resourceKey)

			if len(roles) == 0 {
				return fmt.Errorf("unauthorized") // TODO: do actual json error response here (401)
			}

			allowed, err := a.ACLBiz.AccessAllowed(c.Request().Context(), roles, resource)
			if err != nil {
				return err // TODO: do actual json error response here (500)
			}

			if !allowed {
				return fmt.Errorf("forbidden") // TODO: do actual json error response here (403)
			}

			return next(c)
		}
	}
}

type resourceGetter interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

type accessAllower interface {
	AccessAllowed(ctx context.Context, roles business.Roles, resource interface{}) (bool, error)
}
