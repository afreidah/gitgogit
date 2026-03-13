package auth

import (
	"fmt"

	"gitgogit/config"
)

// SSHProvider implements Provider for SSH key authentication.
// It sets GIT_SSH_COMMAND so git uses the specified private key.
type SSHProvider struct{}

// BuildSSHCommand returns a GIT_SSH_COMMAND value for the given expanded key path.
// The path is shell-quoted to handle spaces.
// StrictHostKeyChecking=accept-new auto-accepts new host keys but rejects changed ones.
// BatchMode=yes prevents SSH from blocking on passphrase prompts.
func BuildSSHCommand(keyPath string) string {
	return fmt.Sprintf("ssh -i %q -o StrictHostKeyChecking=accept-new -o BatchMode=yes", keyPath)
}

func (SSHProvider) Prepare(rawURL string, cfg config.AuthConfig) (string, Env, error) {
	if cfg.Key == "" {
		return rawURL, nil, fmt.Errorf("ssh auth: key path is required")
	}
	// cfg.Key is already ~ expanded by config.Load.
	env := Env{"GIT_SSH_COMMAND=" + BuildSSHCommand(cfg.Key)}
	return rawURL, env, nil
}
