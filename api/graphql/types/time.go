package types

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprint(t.Unix())))
	})
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	timestamp, ok := v.(int64)
	if !ok {
		return time.Time{}, fmt.Errorf("date must be int64")
	}

	return time.Unix(timestamp, 0), nil
}
