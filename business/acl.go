package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/classic-massok/classic-massok-be/lib"
)

// Application Scope types
const (
	appScopeGlobal ApplicationScope = "global"
	appScopeUsers  ApplicationScope = "users"
)

// User Role types (ADD PERM DEFS)
const (
	roleAdmin RoleType = "admin"
	roleUser  RoleType = "user"
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

type ApplicationScope string

func (a ApplicationScope) String() string {
	return string(a)
}

func (a ApplicationScope) Validate() bool {
	switch a {
	case appScopeGlobal:
	case appScopeUsers:
	default:
		// log what was input here
		return false
	}

	return true
}

type RoleType string

func (r RoleType) String() string {
	return string(r)
}

func (r RoleType) Validate() bool {
	switch r {
	case roleAdmin:
	case roleUser:
	default:
		// log what was input here
		return false
	}

	return true
}

// TODO: think out creating an actual role type (we now have a roleType type)... need to ensure we can have ranged (all obj) and
// specific (specific objs) permission sets
// REMEMBER TO VALIDATE ROLES & ROLE TYPES
func GenerateRole(applicationScope ApplicationScope, roleType RoleType, resourceIDs ...string) string {
	if applicationScope == appScopeGlobal {
		return fmt.Sprintf("%s.%s", applicationScope, roleType)
	}

	rIDs := "*"
	if len(resourceIDs) != 0 {
		rIDs = strings.Join(resourceIDs, "|")
	}

	return fmt.Sprintf("%s.%s.%s", applicationScope, roleType, rIDs)
}

type aclRoot interface {
	acl() ACL
}

type aclParent interface {
	parentACL(getter interface{}, id string)
}

type resourceGetter interface {
	Get(ctx context.Context)
}
