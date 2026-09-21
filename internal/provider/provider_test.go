package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveClientAssertionPrefersConfigured(t *testing.T) {
	t.Setenv("ARM_OIDC_TOKEN", "from-env")

	got, err := resolveClientAssertion("configured")
	if err != nil || got != "configured" {
		t.Errorf("expected configured assertion, got %q (err %v)", got, err)
	}
}

func TestResolveClientAssertionFromEnvAndFile(t *testing.T) {
	t.Setenv("ARM_OIDC_TOKEN", "from-env")
	if got, _ := resolveClientAssertion(""); got != "from-env" {
		t.Errorf("expected ARM_OIDC_TOKEN, got %q", got)
	}

	t.Setenv("ARM_OIDC_TOKEN", "")
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte("from-file\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARM_OIDC_TOKEN_FILE_PATH", path)
	if got, _ := resolveClientAssertion(""); got != "from-file" {
		t.Errorf("expected file contents without trailing newline, got %q", got)
	}

	t.Setenv("ARM_OIDC_TOKEN_FILE_PATH", filepath.Join(t.TempDir(), "missing"))
	if _, err := resolveClientAssertion(""); err == nil {
		t.Error("expected an error for a missing token file")
	}
}

func TestResolveClientAssertionEmptyWhenNothingSet(t *testing.T) {
	t.Setenv("ARM_OIDC_TOKEN", "")
	t.Setenv("ARM_OIDC_TOKEN_FILE_PATH", "")
	if got, err := resolveClientAssertion(""); got != "" || err != nil {
		t.Errorf("expected empty assertion, got %q (err %v)", got, err)
	}
}
