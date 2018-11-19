package crond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gamelife1314/crontab/common"
	"github.com/stretchr/testify/assert"
)

func TestHandleJobSave(t *testing.T) {
	initJobManager(t)
	values := url.Values{
		"jobName":  []string{"hello"},
		"command":  []string{"echo hello"},
		"cronExpr": []string{"* * * * *"}}.
		Encode()
	reader := strings.NewReader(values)
	req, err := http.NewRequest(http.MethodPost, "/job/save", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleJobSave)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
}

func TestHandleJobDelete(t *testing.T) {
	initJobManager(t)
	values := url.Values{
		"jobName": []string{"hello"}}.
		Encode()
	reader := strings.NewReader(values)
	req, err := http.NewRequest(http.MethodPost, "/job/delete", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleJobDelete)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
}

func TestHandleJobList(t *testing.T) {
	initJobManager(t)
	req, err := http.NewRequest(http.MethodGet, "/job/list", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleJobList)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
}

func TestHandleJobKill(t *testing.T) {
	initJobManager(t)
	values := url.Values{
		"jobName": []string{"hello"}}.
		Encode()
	reader := strings.NewReader(values)
	req, err := http.NewRequest(http.MethodPost, "/job/kill", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleJobKill)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
	t.Log(resp)
}

func TestHandleJobLog(t *testing.T) {
	initLogManager(t)
	values := url.Values{
		"jobName": []string{"hello"},
		"skip":    []string{"0"},
		"limit":   []string{"1"},
	}.Encode()
	reader := strings.NewReader(values)
	req, err := http.NewRequest(http.MethodPost, "/job/log", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleJobLog)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
	t.Log(resp)
}

func TestHandleWorkerList(t *testing.T) {
	initWorkManager(t)
	req, err := http.NewRequest(http.MethodGet, "/worker/list", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleWorkerList)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	resp := &common.Response{}
	if err := json.Unmarshal(rr.Body.Bytes(), resp); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, resp.ErrorCode)
	t.Log(resp)
}
