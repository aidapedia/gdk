package config

type SecretType string

// SecretType is the type of secret to use.
// We will implement the following secret types:
// - file: read secret from a file
// - gsm: read secret from Google Secret Manager
// - vault: read secret from HashiCorp Vault
const (
	SecretTypeFile  SecretType = "file"
	SecretTypeGSM   SecretType = "gsm"
	SecretTypeVault SecretType = "vault"
)
