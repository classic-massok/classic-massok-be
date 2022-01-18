package business

import (
	"context"
	"strings"

	"github.com/classic-massok/classic-massok-be/business/models"
	"github.com/pkg/errors"
)

func NewACLBiz(verbose bool, resourceRepoBiz resourceRepoBiz) *aclBiz {
	return &aclBiz{resourceRepoBiz, verbose}
}

// aclBiz represents the access control layer business logic
type aclBiz struct {
	resourceRepoBiz resourceRepoBiz
	verbose         bool // TODO: implement verbose logging logic when adding logging
}

func (a *aclBiz) AccessAllowed(ctx context.Context, resource interface{}, action, userID string, roles models.Roles) (bool, error) {
	if resource == nil {
		resource = &RootACL{}
	}

	for {
		hasRootACL, ok := resource.(aclRoot)
		if ok {
			aclItems := hasRootACL.acl()
			for _, aclItem := range aclItems {
				for _, role := range roles {
					if role == models.GlobalAdmin {
						return true, nil
					}

					switch role {
					case models.GlobalAdmin:
						return true, nil
					case models.UserSelf:
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

type resourceRepoBiz interface {
	Get(ctx context.Context, resourceType, resourceID string) (interface{}, error)
}

type aclRoot interface {
	acl() models.ACL
}

type aclParent interface {
	parentACL() (resourceType, id string)
}
