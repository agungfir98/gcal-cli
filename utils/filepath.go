package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var configDir = ".config/gcal-cli"
var Credentials = "credentials.json"
var TokenFile = "token.json"

func CreateConfigPath() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to read home direcotry: %v\n", err)
	}

	dirPath := filepath.Join(homeDir, ".config", "gcal-cli")

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Fatalf("unable to create config file: %v\n", err)
		}
		fmt.Println("config dir created: ", dirPath)
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("directory alredy exist")
	}
}

func GetCredentialsFile() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find homedir: %v\n", err)
	}

	return filepath.Join(homeDir, configDir, Credentials)
}

func GetTokenFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to find homedir: %v\n", err)
	}

	return filepath.Join(homeDir, configDir, TokenFile)

}
