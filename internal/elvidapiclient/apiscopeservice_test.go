package elvidapiclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Minimal coverage of ValidateNewApiScope, the new client function powering
// ApiScopeResource.ModifyPlan (plan-time H2 enforcement). Catches the most
// realistic regression: a typo in the URL path or a wrong status-code branch
// silently disabling plan-time validation. The wider provider test suite is
// deliberately out of scope — see the discussion in elvid PR #833.

func TestValidateNewApiScope_Success_ReturnsNil(t *testing.T) {
	var receivedPath, receivedMethod string
	var receivedBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := ValidateNewApiScope(context.Background(), server.URL, "fake-token", &ApiScopeDto{
		Name:                "convey.api.test",
		AllowMachineClients: true,
	})

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if receivedPath != "/api/ApiScope/validate" {
		t.Errorf("expected POST to /api/ApiScope/validate, got %s", receivedPath)
	}
	if receivedMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", receivedMethod)
	}
	if !strings.Contains(receivedBody, `"convey.api.test"`) {
		t.Errorf("expected DTO JSON in body, got: %s", receivedBody)
	}
}

func TestValidateNewApiScope_BadRequest_ReturnsErrorWithBody(t *testing.T) {
	const elvidMessage = "AD group 'systemaccess-foo-developer' does not exist. Create the group first ... See ADR 2026-05-CORE-2650 (H2)."

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(elvidMessage))
	}))
	defer server.Close()

	err := ValidateNewApiScope(context.Background(), server.URL, "fake-token", &ApiScopeDto{
		Name:                "foo.bar",
		AllowMachineClients: true,
	})

	if err == nil {
		t.Fatal("expected error on 400 response, got nil")
	}
	// The error wraps the response body; Terraform shows this in the plan output.
	if !strings.Contains(err.Error(), elvidMessage) {
		t.Errorf("expected error to contain ElvID's message, got: %v", err)
	}
}
