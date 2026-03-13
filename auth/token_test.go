package auth

import (
	"strings"
	"testing"

	"gitgogit/config"
)

func TestInjectToken_HTTPS(t *testing.T) {
	result, err := InjectToken("https://gitlab.com/org/repo.git", "mytoken")
	if err != nil {
		t.Fatalf("InjectToken() error: %v", err)
	}
	if !strings.Contains(result, "oauth2:mytoken@") {
		t.Errorf("InjectToken() result missing credentials: %q", result)
	}
	if !strings.Contains(result, "gitlab.com") {
		t.Errorf("InjectToken() result missing host: %q", result)
	}
}

func TestInjectToken_HTTP(t *testing.T) {
	result, err := InjectToken("http://example.com/repo.git", "tok")
	if err != nil {
		t.Fatalf("InjectToken() error: %v", err)
	}
	if !strings.Contains(result, "oauth2:tok@") {
		t.Errorf("InjectToken() result missing credentials: %q", result)
	}
}

func TestInjectToken_NonHTTPS(t *testing.T) {
	_, err := InjectToken("git@github.com:org/repo.git", "token")
	if err == nil {
		t.Error("expected error for non-https URL")
	}
}

func TestInjectToken_EmptyToken(t *testing.T) {
	_, err := InjectToken("https://gitlab.com/org/repo.git", "")
	if err == nil {
		t.Error("expected error for empty token")
	}
}

func TestInjectToken_InvalidURL(t *testing.T) {
	_, err := InjectToken("://bad url", "token")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestTokenProvider_Prepare(t *testing.T) {
	t.Setenv("TEST_GITLAB_TOKEN", "secrettoken")

	p := TokenProvider{}
	cfg := config.AuthConfig{Type: "token", Env: "TEST_GITLAB_TOKEN"}
	resolved, env, err := p.Prepare("https://gitlab.com/org/repo.git", cfg)
	if err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}
	if !strings.Contains(resolved, "oauth2:secrettoken@") {
		t.Errorf("Prepare() URL missing credentials: %q", resolved)
	}
	if len(env) != 0 {
		t.Errorf("token auth should not set extra env vars, got %v", env)
	}
}

func TestTokenProvider_MissingEnvVar(t *testing.T) {
	p := TokenProvider{}
	cfg := config.AuthConfig{Type: "token", Env: "DEFINITELY_NOT_SET_XYZ_12345"}
	_, _, err := p.Prepare("https://gitlab.com/org/repo.git", cfg)
	if err == nil {
		t.Error("expected error when env var not set")
	}
}

func TestTokenProvider_MissingEnvName(t *testing.T) {
	p := TokenProvider{}
	cfg := config.AuthConfig{Type: "token", Env: ""}
	_, _, err := p.Prepare("https://gitlab.com/org/repo.git", cfg)
	if err == nil {
		t.Error("expected error when env var name is empty")
	}
}
