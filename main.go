package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	flagDir = "dir"
)

var errDirNotFound = errors.New("could not find script dir")
var errNoArgs = errors.New("script name not provided")
var errCallScript = errors.New("script execution failed")

func main() {
	log.SetFormatter(&log.TextFormatter{})
	app := cli.NewApp()
	app.Name = "qwerty"
	app.Flags = globalFlags()
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run a script",
			Action:  runCommand,
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "scaffold a script.d directory",
			Action:  initCommand,
		},
	}
	app.Run(os.Args)
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  flagDir,
			Value: "",
			Usage: "set the working directory",
		},
	}
}

func getCwd(c *cli.Context) (string, error) {
	cwd := c.String(flagDir)
	if cwd == "" {
		dir, err := os.Executable()
		if err != nil {
			return "", err
		}
		cwd = path.Dir(dir)
	}
	return cwd, nil
}

func execDir(dir, cmd string, things ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c := exec.Command(cmd, things...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	c.Dir = dir
	err := c.Run()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}
