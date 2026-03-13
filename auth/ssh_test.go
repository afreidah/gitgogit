package auth

import (
	"strings"
	"testing"

	"gitgogit/config"
)

func TestBuildSSHCommand(t *testing.T) {
	cmd := BuildSSHCommand("/home/user/.ssh/id_ed25519")
	if !strings.Contains(cmd, "/home/user/.ssh/id_ed25519") {
		t.Errorf("BuildSSHCommand() missing key path: %q", cmd)
	}
	if !strings.Contains(cmd, "BatchMode=yes") {
		t.Errorf("BuildSSHCommand() missing BatchMode=yes: %q", cmd)
	}
	if !strings.Contains(cmd, "StrictHostKeyChecking=accept-new") {
		t.Errorf("BuildSSHCommand() missing StrictHostKeyChecking=accept-new: %q", cmd)
	}
}

func TestBuildSSHCommand_QuotesPath(t *testing.T) {
	cmd := BuildSSHCommand("/path with spaces/id_ed25519")
	// Path with spaces should be quoted
	if !strings.Contains(cmd, `"`) {
		t.Errorf("BuildSSHCommand() path with spaces not quoted: %q", cmd)
	}
}

func TestSSHProvider_Prepare(t *testing.T) {
	p := SSHProvider{}
	cfg := config.AuthConfig{Type: "ssh", Key: "/home/user/.ssh/id_ed25519"}
	url, env, err := p.Prepare("git@github.com:org/repo.git", cfg)
	if err != nil {
		t.Fatalf("Prepare() error: %v", err)
	}
	if url != "git@github.com:org/repo.git" {
		t.Errorf("URL should be unchanged for SSH, got %q", url)
	}
	if len(env) == 0 {
		t.Fatal("expected non-empty env")
	}
	found := false
	for _, e := range env {
		if strings.HasPrefix(e, "GIT_SSH_COMMAND=") {
			found = true
			if !strings.Contains(e, cfg.Key) {
				t.Errorf("GIT_SSH_COMMAND missing key path: %q", e)
			}
		}
	}
	if !found {
		t.Error("GIT_SSH_COMMAND not in env")
	}
}

func TestSSHProvider_MissingKey(t *testing.T) {
	p := SSHProvider{}
	cfg := config.AuthConfig{Type: "ssh", Key: ""}
	_, _, err := p.Prepare("git@github.com:org/repo.git", cfg)
	if err == nil {
		t.Error("expected error when key is empty")
	}
}
