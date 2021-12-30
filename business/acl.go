package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/classic-massok/classic-massok-be/lib"
)

// Resource Refs
const (
	AllResources = "*"
	SelfResource = "self"
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

func (a ApplicationScope) Validate() error {
	switch a {
	case appScopeGlobal:
	case appScopeUsers:
	default:
		// log what was input here
		return fmt.Errorf("invalid application scope provided")
	}

	return nil
}

type RoleType string

func (r RoleType) String() string {
	return string(r)
}

func (r RoleType) Validate() error {
	switch r {
	case roleAdmin:
	case roleUser:
	default:
		// log what was input here
		return fmt.Errorf("invalid role type provided")
	}

	return nil
}

type role string

func (r role) String() string {
	return string(r)
}

func (r role) Validate() error {
	parts := strings.Split(r.String(), ".")
	switch len(parts) {
	case 2:
	case 3:
	default:
		return fmt.Errorf("role is misconfigured")
	}

	// TODO: combine both validations into single error (need multiline error type)
	appScope := ApplicationScope(parts[0])
	if err := appScope.Validate(); err != nil {
		return err
	}

	roleType := RoleType(parts[1])
	if err := roleType.Validate(); err != nil {
		return err
	}

	return nil
}

func generateRoles(applicationScope ApplicationScope, roleType RoleType, resourceIDs ...string) ([]string, error) {
	if applicationScope == appScopeGlobal {
		r := role(fmt.Sprintf("%s.%s", applicationScope, roleType))
		if err := r.Validate(); err != nil {
			return []string{}, fmt.Errorf("invalid role: %w", err)
		}

		return []string{r.String()}, nil
	}

	if len(resourceIDs) == 0 {
		return []string{}, fmt.Errorf("resourceIDs must be provided")
	}

	var roles []string
	for _, resourceID := range resourceIDs {
		r := role(fmt.Sprintf("%s.%s.%s", applicationScope, roleType, resourceID))
		if err := r.Validate(); err != nil {
			return []string{}, fmt.Errorf("invalid role: %w", err)
		}

		roles = append(roles, r.String())
	}

	return roles, nil
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
