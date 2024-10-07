package browser

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

type BrowserOpener interface {
	Open(url string) error
}

func GetBrowserOpener() (BrowserOpener, error) {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return WslOpener{}, nil
	}

	switch runtime.GOOS {
	case "linux":
		return LinuxOpener{}, nil
	case "darwin":
		return MacOpener{}, nil
	case "windows":
		return WindowsOpener{}, nil
	default:
		return nil, errors.New("unsupported platform")

	}
}

type LinuxOpener struct{}

func (o LinuxOpener) Open(url string) error {
	cmd := exec.Command("xdg-open", url)
	return cmd.Start()
}

type WindowsOpener struct{}

func (o WindowsOpener) Open(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	return cmd.Start()
}

type MacOpener struct{}

func (o MacOpener) Open(url string) error {
	cmd := exec.Command("open", url)
	return cmd.Start()
}

type WslOpener struct{}

func (o WslOpener) Open(url string) error {
	cmd := exec.Command("wslview", url)
	return cmd.Start()
}
