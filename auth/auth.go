package auth

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNoUserInToken = errors.New("no user sent in token")
)

type (
	AuthToken struct {
		Token     string `json:"authToken,omitempty" gorethink:"authToken"`
		UserAgent string `json:"userAgent,omitempty" gorethink:"userAgent"`
	}

	AccessToken struct {
		Token    string
		Username string
	}

	ServiceKey struct {
		Key         string `json:"key,omitempty" gorethink:"key"`
		Description string `json:"description,omitempty" gorethink:"description"`
	}

	Authenticator interface {
		Authenticate(username, password, hash string) (bool, error)
		GenerateToken() (string, error)
		IsUpdateSupported() bool
		Name() string
	}
)

func Hash(data string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(h[:]), err
}

func GenerateToken() (string, error) {
	return Hash(time.Now().String())
}

// GetAccessToken returns an AccessToken from the access header
func GetAccessToken(authToken string) (*AccessToken, error) {
	parts := strings.Split(authToken, ":")

	if len(parts) != 2 {
		return nil, ErrNoUserInToken

	}

	return &AccessToken{
		Username: parts[0],
		Token:    parts[1],
	}, nil

}
