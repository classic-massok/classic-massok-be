package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/classic-massok/classic-massok-be/lib"
	"github.com/pkg/errors"
)

// Application Scope types
const (
	appScopeGlobal ApplicationScope = "global"
	appScopeUsers  ApplicationScope = "users"
)

// User Role types (TODO: add permission definitions here)
const (
	roleTypeAdmin RoleType = "admin"
	roleTypeUser  RoleType = "user"
)

// Default roles
const (
	globalAdmin = "global.admin"
	userSelf    = "users.user.self"
)

func NewACLBiz(verbose bool, resourceRepoBiz resourceRepoBiz) *aclBiz {
	return &aclBiz{resourceRepoBiz, verbose}
}

// aclBiz represents the access control layer business logic
type aclBiz struct {
	resourceRepoBiz resourceRepoBiz
	verbose         bool // TODO: implement verbose logging logic when adding logging
}

type ACL []ACE

type ACE struct {
	Roles   Roles
	Actions lib.StringSet
}

func (a *aclBiz) AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles Roles) (bool, error) {
	if resource == nil {
		resource = &RootACL{}
	}

	for {
		hasRootACL, ok := resource.(aclRoot)
		if ok {
			aclItems := hasRootACL.acl()
			for _, aclItem := range aclItems {
				for _, role := range roles {
					if role == globalAdmin {
						return true, nil
					}

					switch role {
					case globalAdmin:
						return true, nil
					case userSelf:
						role = strings.Replace(role, "self", userID, 1)
					}

					if !aclItem.Roles.HasRole(role) {
						continue
					}

					if aclItem.Actions.Contains(action) {
						return true, nil
					}
				}
			}
		}

		hasParentACL, ok := resource.(aclParent)
		if ok {
			resourceType, resourceID := hasParentACL.parentACL()

			var err error
			if resource, err = a.resourceRepoBiz.Get(ctx, resourceType, resourceID); err != nil {
				return false, errors.Wrapf(err, "failed to get %s:%s", resourceType, resourceID)
			}

			continue
		}

		return false, nil
	}
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
		// TODO: log what was input here
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
	case roleTypeAdmin:
	case roleTypeUser:
	default:
		// TODO: log what was input here
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

type Roles []string

func (r Roles) SetRoles(appScope ApplicationScope, roleType RoleType, resourceIDs ...string) error {
	if r == nil {
		r = Roles{}
	}

	roles, err := generateRoles(appScope, roleType, resourceIDs...)
	if err != nil {
		return err
	}

	r = append(r, roles...)
	return nil
}

func (r Roles) HasRole(role string) bool {
	if r == nil {
		return false
	}

	for _, curRole := range r {
		if curRole == role {
			return true
		}
	}

	return false
}

func (r Roles) RemoveRole(role string) bool {
	if r == nil {
		return false
	}

	for i, curRole := range r {
		if curRole == role {
			r = append(r[:i], r[i+1:]...)
			return true
		}
	}

	return false
}

func (r *Roles) DeDupe() { // TODO:  use stringset here
	exists := map[string]struct{}{}
	deduped := Roles{}

	for _, role := range *r {
		if _, ok := exists[role]; !ok {
			exists[role] = struct{}{}
			deduped = append(deduped, role)
		}
	}

	*r = deduped
}

// TODO: do we need this?
func (r Roles) Validate() error {
	if r == nil {
		return nil
	}

	for _, curRole := range r {
		role(curRole).Validate()
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

// ResourceRole is used as a type for a configureable resource id role (i.e. users.user.<userID>)
type ResourceRole string

func (r ResourceRole) String() string {
	return string(r)
}

func (r ResourceRole) Populate(roleType RoleType, resourceID string) string {
	// TODO: should we validate role type here? panic if setup wrong?
	return fmt.Sprintf(r.String(), roleType, resourceID)
}

type resourceRepoBiz interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

type aclRoot interface {
	acl() ACL
}

type aclParent interface {
	parentACL() (resourceType, id string)
}
