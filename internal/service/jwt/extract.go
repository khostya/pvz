package jwt

import (
	"strings"
)

const (
	headerPrefix = "Bearer"
)

func getToken(header string) (string, error) {
	headerSplit := strings.Split(header, " ")
	if len(headerSplit) != 2 || headerSplit[0] != headerPrefix {
		return "", ErrInvalidToken
	}
	return headerSplit[1], nil
}
