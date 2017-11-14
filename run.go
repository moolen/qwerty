package main

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func runCommand(c *cli.Context) error {
	var err error
	var task string
	argv := c.Args()
	if len(argv) == 0 || argv[0] == "" {
		return cli.NewExitError(errNoArgs, 1)
	}
	task = argv[0]
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	foundDir, err := traverseDir(cwd, "script.d")
	if err != nil {
		return cli.NewExitError(err, 3)
	}
	log.Infof("found dir at %s", foundDir)
	stdout, stderr, err := callScript(foundDir, task)
	if err != nil {
		log.Errorln(stdout, stderr)
		return cli.NewExitError(errCallScript, 4)
	}
	return nil
}

func callScript(scriptDir, task string) (stdout string, stderr string, err error) {
	script := path.Join(scriptDir, task)
	execRoot := path.Dir(scriptDir)
	info, err := os.Stat(script)
	if err != nil {
		return "", "", err
	}
	log.Infof("starting at %s", script)
	if info.IsDir() {
		return stdout, stderr, filepath.Walk(script, func(scriptPath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				log.Infof("entering %s", info.Name())
				return nil
			}
			log.Infof("executing traversed file %s", scriptPath)
			stdo, stdr, err := execDir(execRoot, "bash", scriptPath)
			stdout += stdo
			stderr += stdr
			if err != nil {
				return err
			}
			return nil
		})
	}
	log.Infof("executing file %s", script)
	stdout, stderr, err = execDir(execRoot, "bash", script)
	if err != nil {
		return stdout, stderr, err
	}
	return stdout, stderr, err
}

func traverseDir(dir, scriptDir string) (string, error) {
	var err error
	for {
		dir = path.Join(dir, scriptDir)
		log.Infof("trying %s", dir)
		_, err = os.Stat(dir)
		if err == nil || !os.IsNotExist(err) || dir == "/"+scriptDir {
			break
		}
		dir = path.Dir(path.Dir(dir))
	}
	if err != nil {
		return "", errDirNotFound
	}
	return dir, nil
}
