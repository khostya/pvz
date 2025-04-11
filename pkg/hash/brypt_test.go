package hash

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcrypt(t *testing.T) {
	hasher := NewBcryptHash(5)

	password := "12312"
	hashed, err := hasher.Hash(password)
	require.NoError(t, err)
	require.NotEqual(t, hashed, password)

	isEqual := hasher.Equal(EqualsParam{V: password, Hashed: hashed})
	require.True(t, isEqual)

	isEqual = hasher.Equal(EqualsParam{V: password + password, Hashed: hashed})
	require.False(t, isEqual)

	hashed2, err := hasher.Hash(password + password)
	require.NoError(t, err)
	require.NotEqual(t, hashed, hashed2)
}
