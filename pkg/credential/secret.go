package credential

// NewSecretTextCredential creates a secret text credential instance
func NewSecretTextCredential(id, secret string) *StringCredentials {
	return &StringCredentials{
		Credential: Credential{
			Scope:        GLOBALScope,
			ID:           id,
			Class:        SecretTextCredentialStaplerClass,
			StaplerClass: SecretTextCredentialStaplerClass,
		},
		Secret: secret,
	}
}
