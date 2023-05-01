package password

import (
	"glitchz/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := utils.RandomStr(10)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = ComparePassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := utils.RandomStr(11)
	err = ComparePassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	require.Error(t, err)

	_hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, _hashedPassword)
	require.NotEqual(t, hashedPassword, _hashedPassword)
}
