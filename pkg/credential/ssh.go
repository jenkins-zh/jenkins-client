package credential

// NewSSHCredential creates a SSH credential instance
func NewSSHCredential(id, username, passphrase, privateKey string) *SSHCredential {
	keySource := PrivateKeySource{
		StaplerClass: DirectSSHCrenditalStaplerClass,
		PrivateKey:   privateKey,
	}

	return &SSHCredential{
		Credential: Credential{
			Scope:        GLOBALScope,
			ID:           id,
			Class:        SSHCrenditalStaplerClass,
			StaplerClass: SSHCrenditalStaplerClass,
		},
		Username:   username,
		Passphrase: passphrase,
		KeySource:  keySource,
	}
}
