package handler

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestLogPairMap(t *testing.T) {
	lm := logPairMap{}
	lm.Store("Hello", &logPair{})
	val, ok := lm.Load("Hello")

	require.True(t, ok)
	require.NotNil(t, val)

	lm.Delete("Hello")

	val, ok = lm.Load("Hello")
	require.False(t, ok)
	require.Nil(t, val)
}
