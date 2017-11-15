package main

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func muxCommand(c *cli.Context) error {
	var err error
	execList := map[string][]string{}
	argv := c.Args()
	if len(argv) == 0 || argv[0] == "" {
		return cli.NewExitError(errNoArgs, 1)
	}
	cwd, err := getCwd(c)
	if err != nil {
		return cli.NewExitError(err, 2)
	}

	projectDirs, err := ioutil.ReadDir(cwd)
	if err != nil {
		return cli.NewExitError(err, 4)
	}

	for _, projectDir := range projectDirs {
		subProject := path.Join(cwd, projectDir.Name())
		if projectDir.IsDir() {
			log.Debugf("entering %s", projectDir.Name())
			subProjectScriptd := path.Join(subProject, SCRIPTD)
			log.Debugf("checking %s", path.Join(projectDir.Name(), SCRIPTD))
			// ./<subproject>/script.d exists?
			_, err := os.Stat(subProjectScriptd)
			if err != nil && os.IsNotExist(err) {
				continue
			}
			// let's add ./<subproject>/script.d/<task>
			for _, task := range argv {
				if taskExists(subProjectScriptd, task) {
					projectTaskList := execList[task]
					if projectTaskList == nil {
						projectTaskList = []string{}
					}
					log.Debugf("adding task %s in dir %s", task, subProjectScriptd)
					execList[task] = append(projectTaskList, subProjectScriptd)
				}
			}
		}
	}
	if len(execList) == 0 {
		return cli.NewExitError(errNoTasks, 5)
	}
	// @todo(mj): channel/sync.WaitGroup based execution for actual multiplexing
	for task, projectRootDir := range execList {
		log.Infof("running %s", task)
		for _, dir := range projectRootDir {
			log.Infof("  -- %s", path.Base(path.Dir(dir)))
			_, _, err := callScript(dir, task)
			if err != nil {
				return cli.NewExitError(err, 6)
			}
		}
	}
	return nil
}

func taskExists(scriptdPath, task string) bool {
	_, err := os.Stat(path.Join(scriptdPath, task))
	if err != nil {
		return false
	}
	return true
}
