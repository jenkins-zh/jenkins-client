package credential

// NewKubeConfigCredential create a KubeConfig type of credential
func NewKubeConfigCredential(id, kubeconfig string) *KubeConfigCredential {
	credentialSource := KubeconfigSource{
		StaplerClass: DirectKubeconfigCredentialStaperClass,
		Content:      kubeconfig,
	}

	return &KubeConfigCredential{
		Credential: Credential{
			Scope:        GLOBALScope,
			ID:           id,
			Class:        KubeconfigCredentialStaplerClass,
			StaplerClass: KubeconfigCredentialStaplerClass,
		},
		KubeconfigSource: credentialSource,
	}
}
