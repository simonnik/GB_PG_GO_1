package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConfigSuccessEnvProcess(t *testing.T) {
	port := 5432
	host := "172.17.0.2"
	user := "instabank"
	password := "instabank"
	name := "instabank"
	os.Clearenv()
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PORT", strconv.Itoa(port))
	os.Setenv("DB_NAME", name)
	os.Setenv("DB_PASSWORD", password)
	c, err := BuildConfig()
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, host, c.DB.Host)
	assert.Equal(t, port, c.DB.Port)
	assert.Equal(t, user, c.DB.User)
	assert.Equal(t, password, c.DB.Password)
	assert.Equal(t, name, c.DB.Name)
}
