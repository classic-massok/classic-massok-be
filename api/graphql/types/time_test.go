package types

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_MarshalTime(t *testing.T) {
	var bs []byte
	buf := bytes.NewBuffer(bs)

	now := time.Now()
	writer := MarshalTime(now)

	writer.MarshalGQL(buf)
	require.Equal(t, fmt.Sprintf("%d", now.Unix()), buf.String())
}

func Test_UnmarshalTime(t *testing.T) {
	now := time.Now()
	umarshalledTime, err := UnmarshalTime(now.Unix())

	require.NoError(t, err)
	require.WithinDuration(t, now, umarshalledTime, 1*time.Second)
}

func Test_UnmarshalTime_badValue(t *testing.T) {
	umarshalledTime, err := UnmarshalTime(0)
	require.EqualError(t, err, "date must be int64")
	require.Zero(t, umarshalledTime)
}
