package elvidapiclient

import "testing"

func TestBuildTokenRequestUsesSecretWhenNoAssertion(t *testing.T) {
	form := buildTokenRequest("client", "secret", "")

	if form.Get("client_secret") != "secret" {
		t.Errorf("expected client_secret to be sent, got %q", form.Get("client_secret"))
	}
	if form.Get("client_assertion") != "" || form.Get("client_assertion_type") != "" {
		t.Error("expected no client assertion fields when only a secret is given")
	}
	if form.Get("grant_type") != "client_credentials" || form.Get("scope") != elvidApiScope {
		t.Error("grant_type and scope must always be set")
	}
}

func TestBuildTokenRequestPrefersAssertion(t *testing.T) {
	form := buildTokenRequest("client", "secret", "oidc-token")

	if form.Get("client_assertion") != "oidc-token" {
		t.Errorf("expected client_assertion to be sent, got %q", form.Get("client_assertion"))
	}
	if form.Get("client_assertion_type") != "urn:ietf:params:oauth:client-assertion-type:jwt-bearer" {
		t.Errorf("unexpected client_assertion_type %q", form.Get("client_assertion_type"))
	}
	if form.Get("client_secret") != "" {
		t.Error("client_secret must not be sent together with an assertion")
	}
}
