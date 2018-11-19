package crond

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitLogManager(t *testing.T) {
	initLogManager(t)
}

func TestLogManager_Close(t *testing.T) {
	initLogManager(t)
	G_LogManager.Close()
}

func TestLogManager_ListLog(t *testing.T) {
	initLogManager(t)
	if _, err := G_LogManager.ListLog("name", 0, 1); err != nil {
		t.Fatal(err)
	}
}

func initLogManager(t *testing.T) {
	err := LoadConfig("")
	assert.Equal(t, nil, err)
	if err := InitLogManager(); err != nil {
		t.Fatal(err)
	}
}
