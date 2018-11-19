package crond

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitWorkManager(t *testing.T) {
	initWorkManager(t)
}

func TestWorkManager_Close(t *testing.T) {
	initWorkManager(t)
	G_WorkManager.Close()
}

func TestHandleWorkerList2(t *testing.T) {
	initWorkManager(t)
	if _, err := G_WorkManager.ListWorkers(); err != nil {
		t.Fatal(err)
	}
}

func initWorkManager(t *testing.T) {
	err := LoadConfig("")
	assert.Equal(t, nil, err)
	if err := InitWorkManager(); err != nil {
		t.Fatal(err)
	}
}
