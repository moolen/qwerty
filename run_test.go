package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestRunCmd(t *testing.T) {
	dir, _ := ioutil.TempDir("", "qwerty")
	createTestPkg(dir)
	table := []struct {
		args []string
		err  error
	}{
		{
			args: []string{},
			err:  cli.NewExitError(errNoArgs, 1),
		},
		{
			args: []string{"___________________"},
			err:  cli.NewExitError(errDirNotFound, 3),
		},
		{
			args: []string{"--dir", dir, "install"},
			err:  nil,
		},
		{
			args: []string{"--dir", dir, "3tgaer"},
			err:  cli.NewExitError(errCallScript, 4),
		},
	}

	for i, row := range table {
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		for _, flag := range globalFlags() {
			flag.Apply(fs)
		}
		fs.Parse(row.args)
		c := cli.NewContext(cli.NewApp(), fs, nil)
		err := runCommand(c)
		if !reflect.DeepEqual(err, row.err) {
			t.Fatalf("[%d] expected \n%#v\ngot\n%#v", i, row.err, err)

		}
	}

}
func TestTraverseDir(t *testing.T) {

	dir, _ := ioutil.TempDir("", "qwerty")
	os.MkdirAll(path.Join(dir, "script.d", "install"), os.ModePerm)
	os.MkdirAll(path.Join(dir, "foo", "bar", "script.d", "install"), os.ModePerm)

	table := []struct {
		dir       string
		scriptDir string
		found     string
		err       error
	}{
		{
			dir:       dir,
			scriptDir: "script.d",
			found:     path.Join(dir, "script.d"),
			err:       nil,
		},
		{
			dir:       dir,
			scriptDir: "noop.d",
			found:     "",
			err:       errDirNotFound,
		},
		{
			dir:       path.Join(dir, "foo"),
			scriptDir: "script.d",
			found:     path.Join(dir, "script.d"),
			err:       nil,
		},
		{
			dir:       path.Join(dir, "foo", "bar"),
			scriptDir: "script.d",
			found:     path.Join(dir, "foo", "bar", "script.d"),
			err:       nil,
		},
	}

	for i, row := range table {
		found, err := traverseDir(row.dir, row.scriptDir)
		if found != row.found {
			t.Fatalf("[%d] expected\n%#v got\n%#v", i, row.found, found)
		}
		if !reflect.DeepEqual(err, row.err) {
			t.Fatalf("[%d] expected\n%#v got\n%#v", i, row.err, err)
		}
	}

}

func TestCallScript(t *testing.T) {
	dir, _ := ioutil.TempDir("", "qwerty")
	os.MkdirAll(path.Join(dir, "script.d", "install"), os.ModePerm)
	os.MkdirAll(path.Join(dir, "script.d", "errdir"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "install", "00-yolo"), []byte("#!/bin/bash\necho 'success1'; exit 0"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "install", "02-aulu"), []byte("#!/bin/bash\necho 'success2'; exit 0"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "errdir", "02-aulu"), []byte("#!/bin/bash\necho 'errdir' >&2; exit 125"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "error"), []byte("#!/bin/bash\necho 'err1'; exit 1"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "stderror"), []byte("#!/bin/bash\necho 'err1231' >&2; exit 12"), os.ModePerm)
	ioutil.WriteFile(path.Join(dir, "script.d", "success"), []byte("#!/bin/bash\necho 'succ2'; exit 0"), os.ModePerm)

	table := []struct {
		dir    string
		task   string
		stdout string
		stderr string
		err    string
	}{
		{
			dir:    path.Join(dir, "script.d"),
			task:   "success",
			stdout: "succ2\n",
			stderr: "",
			err:    "",
		},
		{
			dir:    path.Join(dir, "script.d"),
			task:   "install",
			stdout: "success1\nsuccess2\n",
			stderr: "",
			err:    "",
		},
		{
			dir:    path.Join(dir, "script.d"),
			task:   "error",
			stdout: "err1\n",
			stderr: "",
			err:    "exit status 1",
		},
		{
			dir:    path.Join(dir, "script.d"),
			task:   "errdir",
			stdout: "",
			stderr: "errdir\n",
			err:    "exit status 125",
		},
		{
			dir:    path.Join(dir, "script.d"),
			task:   "stderror",
			stdout: "",
			stderr: "err1231\n",
			err:    "exit status 12",
		},
		{
			dir:    dir,
			task:   "noop.d",
			stdout: "",
			stderr: "",
			err:    "stat " + path.Join(dir, "noop.d") + ": no such file or directory",
		},
	}
	for i, row := range table {
		stdout, stderr, err := callScript(row.dir, row.task)
		if stdout != row.stdout {
			t.Fatalf("[%d] expected\n%#v got\n%#v", i, row.stdout, stdout)
		}
		if stderr != row.stderr {
			t.Fatalf("[%d] expected\n%#v got\n%#v", i, row.stderr, stderr)
		}
		if err != nil && err.Error() != row.err {
			t.Fatalf("[%d] expected\n%#v got\n%#v", i, row.err, err.Error())
		}
	}
}
