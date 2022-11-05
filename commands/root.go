package commands

import (
	"docker-symfony/util"
	"fmt"
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
	"os"
)

func CommonCommands() []*console.Command {
	d4dCommands := []*console.Command{
		startCmd,
		passwordShowCmd,
		selfUpdateCmd,
	}

	return d4dCommands
}

func InitAppFunc(c *console.Context) error {

	return nil
}

// WelcomeAction displays a message when no command
func WelcomeAction(c *console.Context) error {
	util.LoadEnvFile(envFile)

	terminal.Println("")
	console.ShowVersion(c)
	terminal.Println(c.App.Usage)
	terminal.Println("")
	terminal.Println("The following information has been set:")
	terminal.Println("")
	terminal.Println("Server IP: 127.0.0.1")
	terminal.Println(fmt.Sprintf("Server Hostname: %s", os.Getenv("PROJECT_DOMAIN_1")))
	terminal.Println("")
	terminal.Println("To login now, follow this link:")
	terminal.Println("")
	terminal.Println(fmt.Sprintf("Project URL: http://%s", os.Getenv("PROJECT_DOMAIN_1")))
	terminal.Println(fmt.Sprintf("phpMyAdmin: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_PMA")))
	terminal.Println("")
	terminal.Println("Extra features:")
	terminal.Println(fmt.Sprintf("MailHog: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_MAILHOG_HTTP")))
	terminal.Println(fmt.Sprintf("RabbitMQ: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_RABBITMQ_MANAGEMENT")))
	terminal.Println(fmt.Sprintf("Elasticsearch: http://%s:%s", os.Getenv("PROJECT_DOMAIN_1"), os.Getenv("PORT_ELASTICSEARCH_HEAD")))
	terminal.Println("")
	terminal.Println("Thank you for using Docker for Symfony. Should you have any questions, don't hesitate to contact us at support@d4d.lt")
	displayCommandsHelp(c, []*console.Command{})
	terminal.Println("")
	terminal.Printf("Show all commands with <info>%s help</>,\n", c.App.HelpName)
	terminal.Printf("Get help for a specific command with <info>%s help COMMAND</>.\n", c.App.HelpName)

	return nil
}

func displayCommandsHelp(c *console.Context, cmds []*console.Command) {
	console.HelpPrinter(c.App.Writer, `{{range .}}  <info>{{.PreferredName}}</>{{"\t"}}{{.Usage}}{{"\n"}}{{end}}`, cmds)
}
