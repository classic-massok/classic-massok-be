package rest

import "github.com/labstack/echo/v4"

func Configure(rest *echo.Group) {
	rest.GET("/hello-world", helloWorld)
}

func helloWorld(c echo.Context) error {
	return c.JSON(200, "Hello, World!")
}
