package config

import (
	"os"
	"path"

	"github.com/joho/godotenv"
)

const CredsDirectoryEnv = "MYCAL_CREDENTIALS_DIRECTORY"

var credsDirectory string

func init() {
	godotenv.Load()
	credsDirectory = os.Getenv(CredsDirectoryEnv)
}

// GetCredsDirectory returns the configured credentials directory
func GetCredsDirectory() string {
	return credsDirectory
}

// GetCredentialsPath returns the full path to the credentials file
func GetCredentialsPath() string {
	return path.Join(credsDirectory, "myCalAppCredentials.json")
}

// GetTokenPath returns the full path to the token file
func GetTokenPath() string {
	return path.Join(credsDirectory, "myCalAppToken.json")
}
