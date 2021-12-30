package authz

import "github.com/labstack/echo/v4"

func LoadResource(c echo.Context, bizRepo resourceGetter, resourceType, resourceID string) error {
	resource, err := bizRepo.Get(c.Request().Context(), resourceType, resourceID)
	if err != nil {
		return err
	}

	c.Set("resource", resource)
	return nil
}
