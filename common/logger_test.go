package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetLogFile(t *testing.T) {
	var (
		err error
		pwd string
	)
	err = SetLogFile("")
	assert.Equal(t, err, nil)

	if pwd, err = os.Getwd(); err == nil {
		assert.FileExists(t, pwd+"/"+"crond.log")
	} else {
		t.Fatal("Some unexpected errors happened.")
	}
	err = SetLogFile("this-is-a-mot-exits-file.log")
	assert.Equal(t, err, nil)
}
