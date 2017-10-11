package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigWithOverrides(t *testing.T) {
	original := Configuration{}
	original.JWT.Secret = "jwt-secret"
	original.DB.Name = "db-name"
	original.DB.ConnURL = "conn-url"
	original.API.Host = "api-host"
	original.API.Port = 12356

	tmpfile, err := ioutil.TempFile("", "gs-test")
	assert.Nil(t, err)

	fname := tmpfile.Name() + ".json"
	err = os.Rename(tmpfile.Name(), fname)
	assert.Nil(t, err)
	defer os.Remove(fname)

	content, err := json.Marshal(&original)
	assert.Nil(t, err)

	err = ioutil.WriteFile(fname, content, 0755)
	assert.Nil(t, err)

	// override some values
	os.Setenv("GS_GEO_JWT_SECRET", "env-jwt-secret")
	os.Setenv("GS_GEO_DB_NAME", "env-db-name")
	os.Setenv("GS_GEO_API_PORT", "456456")

	config, err := Load(fname)
	assert.Nil(t, err)
	assert.NotNil(t, config)

	// check we loaded from the file
	assert.Equal(t, config.DB.ConnURL, original.DB.ConnURL)
	assert.Equal(t, config.API.Host, original.API.Host)

	// check we got the overrides
	assert.Equal(t, "env-jwt-secret", config.JWT.Secret)
	assert.Equal(t, "env-db-name", config.DB.Name)
	assert.EqualValues(t, 456456, config.API.Port)
}