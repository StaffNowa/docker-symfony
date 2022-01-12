package commands

import (
	"docker-symfony/util"
	"fmt"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
	"os"
)

var passwordShowCmd = &console.Command{
	Category: "d4d",
	Name:     "passwd:show",
	Usage:    "",
	Action: func(c *console.Context) error {
		doPasswordShow()

		return nil
	},
}

func doPasswordShow() {
	util.ExecCommand("clear")

	util.LoadEnvFile(envFile)

	terminal.Println("")
	terminal.Println("The following information has been set:")
	terminal.Println("")
	terminal.Println("Server IP: 127.0.0.1")
	terminal.Println(fmt.Sprintf("Server Hostname: %s", os.Getenv("PROJECT_DOMAIN_1")))
	terminal.Println("")
	terminal.Println("MySQL root username: root")
	terminal.Println(fmt.Sprintf("MySQL root password: %s", os.Getenv("MYSQL_ROOT_PASSWORD")))
	terminal.Println("")
	terminal.Println(fmt.Sprintf("MySQL database name: %s", os.Getenv("MYSQL_DATABASE")))
	terminal.Println(fmt.Sprintf("MySQL username: %s", os.Getenv("MYSQL_USER")))
	terminal.Println(fmt.Sprintf("MySQL password: %s", os.Getenv("MYSQL_PASSWORD")))
	terminal.Println("")
	terminal.Println("To login now, follow this link:")
	terminal.Println("")
	terminal.Println(fmt.Sprintf("Project URL: http://%s", os.Getenv("PROJECT_DOMAIN_1")))
	terminal.Println(fmt.Sprintf("phpMyAdmin: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_PMA")))
	terminal.Println(fmt.Sprintf("MailHog: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_MAILHOG_HTTP")))
	terminal.Println("")
}
