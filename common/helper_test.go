package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildResponse(t *testing.T) {
	var (
		resp []byte
		err  error
		data map[string]interface{}
	)

	data = map[string]interface{}{
		"id": 1,
	}

	resp, err = BuildResponse(0, "success", data)

	assert.Equal(t, err, nil)
	if len(resp) == 0 {
		t.Fatal("Errors happened when common.BuildResponse build response.")
	}
}

func TestExtractWorkerIp(t *testing.T) {
	key := "/cron/worker/192.168.0.1"
	assert.Equal(t, "192.168.0.1", ExtractWorkerIp(key))
}
