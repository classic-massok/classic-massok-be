package api

import (
	"net/http"

	"github.com/classic-massok/classic-massok-be/api/graphql"
	"github.com/classic-massok/classic-massok-be/api/rest"
	"github.com/labstack/echo/v4"
)

func GetRouter() http.Handler {
	e := echo.New()

	apiRouter := e.Group("/api")

	restRouter := apiRouter.Group("/rest")
	rest.Configure(restRouter)

	graphqlRouter := apiRouter.Group("/graphql")
	graphql.Configure(graphqlRouter)

	return e
}
