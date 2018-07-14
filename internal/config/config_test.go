package config

import (
	"os"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/fs"
)

func TestConfigDefaults(t *testing.T) {
	f := fs.NewFile(t, "config-defaults")
	defer f.Remove()

	configPath := f.Path()
	c, err := Parse("", []string{configPath})
	assert.NilError(t, err)

	assert.Equal(t, c.WarehouseDB.DatabaseName, "events")
}

func TestConfigWillReadFromEnv(t *testing.T) {
	os.Setenv("", "test:1010")
	os.Setenv("COCOAPODS__WAREHOUSE__USERNAME", "test")

	c, err := Parse("", []string{})
	assert.NilError(t, err)

	assert.Equal(t, c.WarehouseDB.Username, "test")
}
