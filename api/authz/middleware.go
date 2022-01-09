package authz

import (
	"context"
	"fmt"

	"github.com/classic-massok/classic-massok-be/api/core"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/labstack/echo/v4"
)

const resourceKey = "Resource"

type AuthzMW struct {
	ACLBiz          accessAllower
	ResourceRepoBiz resourceGetter
}

func (a *AuthzMW) LoadResource(resourceType string, resourceID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := LoadResource(c, a.ResourceRepoBiz, resourceType, resourceID); err != nil {
				return core.JSON(c, 404, nil, err)
			}

			return next(c)
		}
	}
}

func (a *AuthzMW) RequiresPermission(action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			switch err := RequiresPermission(c, a.ACLBiz, action); err {
			case nil:
			case lib.ErrUnauthorized:
				return core.JSON(c, 401, nil, err)
			case lib.ErrForbidden:
				return core.JSON(c, 403, nil, err)
			default:
				return core.JSON(c, 500, nil, fmt.Errorf("server error: %w", err))
			}

			return next(c)
		}
	}
}

type resourceGetter interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

type accessAllower interface {
	AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles business.Roles) (bool, error)
}
