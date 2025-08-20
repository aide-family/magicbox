package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aide-family/magicbox/password"
)

func TestNew_DefaultSalt(t *testing.T) {
	pwd := "mySecretPassword"
	passwordObj, err := password.New(pwd)
	assert.Nil(t, err)
	assert.NotNil(t, passwordObj)
	assert.NotEmpty(t, passwordObj.Salt())
	assert.True(t, passwordObj.Equal(passwordObj.Value()))
	assert.False(t, passwordObj.Equal(""))
}
