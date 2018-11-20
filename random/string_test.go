package random_test

import (
	"testing"

	"github.com/bakku/easyalert/random"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	str, err := random.String(32)
	require.Nil(t, err)

	require.Equal(t, 32, len(str))
}
