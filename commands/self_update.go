package commands

import (
	"docker-symfony/util"
	"fmt"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
	"runtime"
)

var selfUpdateCmd = &console.Command{
	Category: "d4d",
	Name:     "self-update",
	Aliases:  []*console.Alias{{Name: "self-update"}},
	Usage:    "",
	Action: func(c *console.Context) error {
		os := runtime.GOOS
		arch := runtime.GOARCH

		util.ExecCommand("clear")
		terminal.Println("We will process upgrading Docker for Symfony to the latest version.")

		var filename = ""
		if os != "darwin" {
			filename = fmt.Sprintf("d4d_linux_%s.tar.gz", arch)
		} else {
			filename = "d4d_darwin_all.tar.gz"
		}

		util.DownloadFile(filename, "https://github.com/StaffNowa/docker-symfony/releases/latest/download/d4d_darwin_all.tar.gz")
		util.ExecCommand("tar xzfv " + filename)
		terminal.Println("Docker for Symfony successfully upgraded to the latest version.")

		return nil
	},
}
