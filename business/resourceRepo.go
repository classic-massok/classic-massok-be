package business

import (
	"context"
	"fmt"
)

// Resource types
const (
	usersResource = "users"
)

type ResourceRepo struct {
	User userGetter
}

func (r *ResourceRepo) Get(ctx context.Context, resourceType, resourceID string) (interface{}, error) {
	switch resourceType {
	case usersResource:
		return r.User.Get(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}
