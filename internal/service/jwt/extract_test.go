package jwt

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGetToken(t *testing.T) {
	genHeader := func(token string) string {
		return "Bearer " + token
	}

	t.Run("ok", func(t *testing.T) {
		token := "Gesg"
		actual, err := getToken(genHeader(token))

		require.NoError(t, err)
		require.Equal(t, token, actual)
	})

	t.Run("error prefix", func(t *testing.T) {
		token := "B Gesg"
		_, err := getToken(token)

		require.Error(t, err)
	})

	t.Run("empty token", func(t *testing.T) {
		token := ""
		_, err := getToken(token)
		require.Error(t, err)
	})

	t.Run("many words in token", func(t *testing.T) {
		token := strings.Repeat("131", 2)
		_, err := getToken(token)
		require.Error(t, err)
	})
}
