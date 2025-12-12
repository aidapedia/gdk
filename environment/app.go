package environment

import (
	"os"
)

func GetAppEnvironment() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	return Development
}

func GetSecretGSMProjectID() string {
	if projectID := os.Getenv("SECRET_GSM_PROJECT_ID"); projectID != "" {
		return projectID
	}
	return ""
}

func GetSecretFilePath() string {
	if path := os.Getenv("SECRET_FILE_PATH"); path != "" {
		return path
	}
	return ""
}

func GetConfigPath() string {
	if path := os.Getenv("CONFIG_FILE_PATH"); path != "" {
		return path
	}
	return ""
}

func GetSecretVaultAddress() string {
	if address := os.Getenv("SECRET_VAULT_ADDRESS"); address != "" {
		return address
	}
	return ""
}

func GetSecretVaultToken() string {
	if token := os.Getenv("SECRET_VAULT_TOKEN"); token != "" {
		return token
	}
	return ""
}

func GetSecretVaultPath() string {
	if path := os.Getenv("SECRET_VAULT_PATH"); path != "" {
		return path
	}
	return ""
}
