package easyalert_test

import (
	"testing"

	"github.com/bakku/easyalert"
	"github.com/stretchr/testify/require"
)

func TestHumanStatus(t *testing.T) {
	var (
		pending = easyalert.Alert{Status: 0}
		sent    = easyalert.Alert{Status: 1}
		failed  = easyalert.Alert{Status: 2}
		invalid = easyalert.Alert{Status: 3}
	)

	require.Equal(t, "pending", pending.HumanStatus())
	require.Equal(t, "sent", sent.HumanStatus())
	require.Equal(t, "failed", failed.HumanStatus())
	require.Equal(t, "invalid status", invalid.HumanStatus())
}
