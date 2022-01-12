package main

import (
	"docker-symfony/commands"
	"docker-symfony/util"
	"fmt"
	"github.com/symfony-cli/console"
	"os"
	"time"
)

var (
	version   = "dev"
	channel   = "dev"
	buildDate string
)

func Init() {
	currentDir := util.GetCurrentDir()

	// Added for security
	util.Chmod(currentDir+"/docker", 0700)
	util.Chmod(currentDir, 0700)

	env := currentDir + "/.env"
	envDist := env + ".dist"

	envSecret := env + ".secret"
	envSecretDist := envSecret + ".dist"

	if !util.FileExists(env) {
		util.Copy(envDist, env)
	}

	if !util.FileExists(envSecret) {
		util.Copy(envSecretDist, envSecret)
	}
}

func main() {
	Init()

	args := os.Args

	cmds := commands.CommonCommands()
	app := &console.Application{
		Name:      "D4D",
		Usage:     "Docker Symfony gives you everything you need to develop a Symfony application. This complete stack runs with docker and docker-compose.",
		Copyright: fmt.Sprintf("(c) 2018-%d D4D", time.Now().Year()),
		Commands:  cmds,
		Action: func(ctx *console.Context) error {
			if ctx.Args().Len() == 0 {
				return commands.WelcomeAction(ctx)
			}
			return console.ShowAppHelpAction(ctx)
		},
		Before:    commands.InitAppFunc,
		Version:   version,
		Channel:   channel,
		BuildDate: buildDate,
	}
	app.Run(args)
}
