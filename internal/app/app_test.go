package app

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/ut"
)

func TestNewBuildsApp(t *testing.T) {
	openLucky, err := New(Options{ConfigPath: t.TempDir() + "/openlucky.json", AdminUsername: "admin", AdminPassword: "secret"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if openLucky == nil || openLucky.Engine == nil {
		t.Fatal("New returned nil app or engine")
	}
}

func TestAuthAndModuleAPI(t *testing.T) {
	openLucky, err := New(Options{ConfigPath: t.TempDir() + "/openlucky.json", AdminUsername: "admin", AdminPassword: "secret"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	unauthorized := ut.PerformRequest(openLucky.Engine.Engine, "GET", "/api/modules/list", nil)
	if unauthorized.Result().StatusCode() != 401 {
		t.Fatalf("unauthorized status = %d, want 401", unauthorized.Result().StatusCode())
	}

	loginBody := []byte(`{"username":"admin","password":"secret"}`)
	login := ut.PerformRequest(openLucky.Engine.Engine, "POST", "/api/login", &ut.Body{Body: bytes.NewBuffer(loginBody), Len: len(loginBody)}, ut.Header{Key: "Content-Type", Value: "application/json"})
	if login.Result().StatusCode() != 200 {
		t.Fatalf("login status = %d body=%s", login.Result().StatusCode(), string(login.Result().Body()))
	}

	var loginEnvelope struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(login.Result().Body(), &loginEnvelope); err != nil {
		t.Fatalf("login response JSON: %v", err)
	}
	if loginEnvelope.Data.Token == "" {
		t.Fatal("login token is empty")
	}

	modules := ut.PerformRequest(openLucky.Engine.Engine, "GET", "/api/modules/list", nil, ut.Header{Key: "OpenLucky-Admin-Token", Value: loginEnvelope.Data.Token})
	if modules.Result().StatusCode() != 200 {
		t.Fatalf("modules status = %d body=%s", modules.Result().StatusCode(), string(modules.Result().Body()))
	}

	stub := ut.PerformRequest(openLucky.Engine.Engine, "GET", "/api/docker/containers", nil, ut.Header{Key: "OpenLucky-Admin-Token", Value: loginEnvelope.Data.Token})
	if stub.Result().StatusCode() != 501 {
		t.Fatalf("stub status = %d body=%s", stub.Result().StatusCode(), string(stub.Result().Body()))
	}
}
