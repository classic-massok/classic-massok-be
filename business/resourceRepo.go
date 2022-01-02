package business

import (
	"context"
	"fmt"
)

// Resource types
const (
	UsersResource = "users"
)

type ResourceRepo struct {
	User userGetter
}

func (r *ResourceRepo) Get(ctx context.Context, resourceType, resourceID string) (interface{}, error) {
	switch resourceType {
	case UsersResource:
		return r.User.Get(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}
