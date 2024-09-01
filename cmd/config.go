/*
Copyright © 2024 Agung Firmansyah agungfir98@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	pathutils "github.com/agungfir98/gcal-cli/utils/path_utils"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "config",
	Short: "create a .config/ directory at user home directory to store credentials and token json",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seting up gcal-cli")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("unable to read home direcotry: %v\n", err)
		}

		dirPath := filepath.Join(homeDir, pathutils.ConfigDir)

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
			return
		}
		fmt.Println("config directory setup done.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
