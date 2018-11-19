package crond

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	var (
		err          error
		pwd          string
		oldConfigEnv string
	)
	pwd, err = os.Getwd()
	oldConfigEnv = os.Getenv("CROND_CONFIG_FILE")
	os.Unsetenv("CROND_CONFIG_FILE")
	err = LoadConfig("")
	assert.NotEqual(t, err, nil)
	assert.Contains(t, err.Error(), "no such file or directory")
	if oldConfigEnv != "" {
		os.Setenv("CROND_CONFIG_FILE", oldConfigEnv)
	} else {
		var parentDir = filepath.Dir(pwd)
		os.Setenv("CROND_CONFIG_FILE", parentDir+string(os.PathSeparator)+".travis.crond.yml")
	}
	err = LoadConfig("")
	assert.Equal(t, err, nil)
}
