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
	SCRIPTD   = "script.d"
	flagDir   = "dir"
	flagDebug = "debug"
)

var errDirNotFound = errors.New("could not find script dir")
var errNoArgs = errors.New("script name not provided")
var errCallScript = errors.New("script execution failed")
var errNoScript = errors.New("script does not exist")
var errNoTasks = errors.New("no tasks found, nothing to do")

func main() {
	log.SetFormatter(&log.TextFormatter{})
	app := cli.NewApp()
	app.Name = "qwerty"
	app.Flags = globalFlags()
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "runs the specified command, you may specify multiple commands.",
			Action:  runCommand,
		},
		{
			Name:    "mux",
			Aliases: []string{"m"},
			Usage:   "multiplex commands in subdirectories",
			// todo: add flags for --ignore-not-existing tasks
			Action: muxCommand,
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "scaffolds a script.d directory",
			Action:  initCommand,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool(flagDebug) {
			log.SetLevel(log.DebugLevel)
		}
		return nil
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
		cli.BoolFlag{
			Name:  flagDebug,
			Usage: "enable debug",
		},
	}
}

func getCwd(c *cli.Context) (string, error) {
	var execDir string
	cwd := c.GlobalString(flagDir)
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	execDir = path.Dir(dir)
	if !path.IsAbs(cwd) {
		cwd = path.Join(execDir, cwd)
	}

	log.Debugf("cwd: %s", cwd)
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
