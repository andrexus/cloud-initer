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
	original.DB.Path = "db.bolt"
	original.API.Host = "api-host"
	original.API.Port = 8000

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
	os.Setenv("CLOUD_INITER_DB_PATH", "db.bolt")
	os.Setenv("CLOUD_INITER_API_HOST", "api-host")
	os.Setenv("CLOUD_INITER_API_PORT", "8000")

	config, err := Load(fname)
	assert.Nil(t, err)
	assert.NotNil(t, config)

	// check we loaded from the file
	assert.Equal(t, config.DB.Path, original.DB.Path)
	assert.Equal(t, config.API.Host, original.API.Host)

	// check we got the overrides
	assert.Equal(t, "db.bolt", config.DB.Path)
	assert.EqualValues(t, "api-host", config.API.Host)
	assert.EqualValues(t, 8000, config.API.Port)
}