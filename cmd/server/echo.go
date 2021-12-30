package main

import (
	"net/http"

	"github.com/classic-massok/classic-massok-be/api"
	"github.com/classic-massok/classic-massok-be/business"
)

func getEchoRouter() http.Handler { // TODO: figure out if setup of resource repo seem reasonable
	userBiz := business.NewUsersBiz()

	resourceRepo := &business.ResourceRepo{
		userBiz,
	}

	aclBiz := business.NewACLBiz(true, resourceRepo) // TODO: make this configureable

	router := &api.Router{
		aclBiz,
		userBiz,
	}

	return router.SetRouter(resourceRepo)
}
