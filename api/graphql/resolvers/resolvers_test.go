package resolvers

import (
	"testing"

	"github.com/classic-massok/classic-massok-be/api/graphql/resolvers/resolversfakes"
	"github.com/stretchr/testify/require"
)

func TestResolver_Mutation(t *testing.T) {
	r := &Resolver{
		&resolversfakes.FakeUsersBiz{},
	}

	m := r.Mutation()
	require.NotNil(t, m)
}

func TestResolver_Query(t *testing.T) {
	r := &Resolver{
		&resolversfakes.FakeUsersBiz{},
	}

	q := r.Query()
	require.NotNil(t, q)
}
