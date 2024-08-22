package utils

import (
	"fmt"
	"time"
)

func ShowLoading(done chan bool) {
	loadingChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\rFetching... %s", loadingChars[i%len(loadingChars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
