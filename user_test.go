package easyalert_test

import (
	"testing"

	"github.com/bakku/easyalert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	u := easyalert.User{}
	err := u.HashPassword("test1234")

	require.Nil(t, err)
	require.NotEqual(t, "", u.PasswordDigest)
}

func TestValidPassword_ReturnsTrueIfValid(t *testing.T) {
	// equals to test1234
	u := easyalert.User{
		PasswordDigest: "$2a$10$tqsaBkRHWCUnome4ybT3ouoSysxu1UAttDQU0jFv19Qps/Qr8FKmW",
	}

	require.True(t, u.ValidPassword("test1234"))
}

func TestValidPassword_ReturnsFalseIfInvalid(t *testing.T) {
	// equals to test1234
	u := easyalert.User{
		PasswordDigest: "$2a$10$tqsaBkRHWCUnome4ybT3ouoSysxu1UAttDQU0jFv19Qps/Qr8FKmW",
	}

	require.False(t, u.ValidPassword("test123"))
}
