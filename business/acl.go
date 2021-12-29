package business

import (
	"context"

	"github.com/classic-massok/classic-massok-be/lib"
)

func NewACLBiz(verbose bool) *aclBiz {
	return &aclBiz{verbose}
}

// aclBiz represents the access control layer business logic
type aclBiz struct {
	Verbose bool // TODO: implement verbose logging logic when adding logging
}

type ACL []ACE

type ACE struct {
	Roles   Roles
	Actions lib.StringSet
}

func (a *aclBiz) AccessAllowed(ctx context.Context, roles Roles, resource interface{}) (bool, error) {
	if resource == nil {
		resource = &RootACL{}
	}

	for {
		// hasRootACL, ok := resource.(aclRoot)
		// if ok {

		// }

		// hasParentACL, ok := resource.(aclParent)
		// if ok {

		// 	continue
		// }

		return false, nil
	}

	return true, nil
}

type RootACL struct{}

func (r *RootACL) acl() {

}

type aclRoot interface {
	acl() ACL
}

type aclParent interface {
	parent(getter interface{}, id string)
}

type resourceGetter interface {
	Get(ctx context.Context)
}
