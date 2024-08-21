/*
Copyright © 2024 Agung Firmansyah agungfir98@gmail.com
*/
package main

import (
	"gcal-cli/cmd"
	"gcal-cli/utils"
)

func main() {
	utils.CreateConfigPath()
	cmd.Execute()
}
