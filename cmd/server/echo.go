package main

import (
	"net/http"

	"github.com/classic-massok/classic-massok-be/api"
	"github.com/classic-massok/classic-massok-be/business"
	"github.com/classic-massok/classic-massok-be/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func getEchoRouter(db *mongo.Database, cfg *config.Config) http.Handler { // TODO: figure out if setup of resource repo seem reasonable
	userBiz := business.NewUsersBiz(db) // TODO: inject db connection through biz layers, figure out if we want to interface this

	resourceRepo := &business.ResourceRepo{
		userBiz,
	}

	aclBiz := business.NewACLBiz(true, resourceRepo) // TODO: make this configureable

	router := &api.Router{
		aclBiz,
		userBiz,
	}

	return router.SetRouter(resourceRepo, cfg)
}
