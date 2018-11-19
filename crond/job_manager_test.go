package crond

import (
	"github.com/gamelife1314/crontab/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitJobManager(t *testing.T) {
	initJobManager(t)
}

func TestJobManager_SaveJob(t *testing.T) {
	initJobManager(t)
	job := &common.Job{
		Name:     "hello",
		Command:  "echo hello",
		CronExpr: "* * * * *",
	}
	if _, err := G_JobManager.SaveJob(job); err != nil {
		t.Fatal(err)
	}
}

func TestJobManager_DeleteJob(t *testing.T) {
	initJobManager(t)
	if _, err := G_JobManager.DeleteJob("hello"); err != nil {
		t.Fatal(err)
	}
}

func TestJobManager_ListJobs(t *testing.T) {
	initJobManager(t)
	if _, err := G_JobManager.ListJobs(); err != nil {
		t.Fatal(err)
	}
}

func TestJobManager_KillJob(t *testing.T) {
	initJobManager(t)
	if err := G_JobManager.KillJob("hello"); err != nil {
		t.Fatal(err)
	}
}

func initJobManager(t *testing.T) {
	err := LoadConfig("")
	assert.Equal(t, nil, err)
	if err := InitJobManager(); err != nil {
		t.Fatal(err)
	}
}
