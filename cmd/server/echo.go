package main

import (
	"net/http"

	"github.com/classic-massok/classic-massok-be/api"
	"github.com/classic-massok/classic-massok-be/business"
)

func getEchoRouter() http.Handler {
	router := &api.Router{
		business.NewACLBiz(true), // TODO: make this configureable
		business.NewUsersBiz(),
	}

	return router.SetRouter()
}
