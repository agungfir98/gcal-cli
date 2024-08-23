package pathutils

import (
	"log"
	"os"
	"path/filepath"
)

var ConfigDir = ".config/gcal-cli"
var Credentials = "credentials.json"
var TokenFile = "token.json"

func GetConfigPath() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to read home directory: %v\n", err)
	}
	dirPath := filepath.Join(homeDir, ConfigDir)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file not found, initiate gcal-cli first")
	}
}

func GetCredentialsFile() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find homedir: %v\n", err)
	}

	return filepath.Join(homeDir, ConfigDir, Credentials)
}

func GetTokenFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find homedir: %v\n", err)
	}

	return filepath.Join(homeDir, ConfigDir, TokenFile)

}
