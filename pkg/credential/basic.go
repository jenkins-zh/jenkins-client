package credential

// NewUsernamePasswordCredential creates a username password type of credential
func NewUsernamePasswordCredential(id, username, password string) *UsernamePasswordCredential {
	return &UsernamePasswordCredential{
		Credential: Credential{
			Scope:        GLOBALScope,
			ID:           id,
			Class:        UsernamePassswordCredentialStaplerClass,
			StaplerClass: UsernamePassswordCredentialStaplerClass,
		},
		Username: username,
		Password: password,
	}
}
