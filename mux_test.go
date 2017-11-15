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

func TestMuxCommand(t *testing.T) {

	// dir is the root of the pkg
	dir, _ := ioutil.TempDir("", "trollolol")
	createTestPkg(path.Join(dir, "proj1"))
	createTestPkg(path.Join(dir, "proj2"))
	createTestPkg(path.Join(dir, "proj3"))
	createTestPkg(path.Join(dir, "proj4"))

	table := []struct {
		err  error
		code int
		args []string
	}{
		{
			code: 1,
			err:  errNoArgs,
			args: []string{},
		},
		{
			err:  nil,
			args: []string{"--dir", dir, "install"},
		},
		{
			code: 4,
			args: []string{"--dir", "/root", "errdir"},
		},
		{

			code: 5,
			args: []string{"--dir", path.Join(dir, "proj1"), "errdir"},
		},
		{

			code: 6,
			args: []string{"--dir", dir, "errdir"},
		},
	}

	for i, row := range table {
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		for _, flag := range globalFlags() {
			flag.Apply(fs)
		}
		fs.Parse(row.args)
		c := cli.NewContext(cli.NewApp(), fs, nil)
		err := muxCommand(c)
		exitErr, _ := err.(*cli.ExitError)
		if !reflect.DeepEqual(err, row.err) && exitErr.ExitCode() != row.code {
			t.Fatalf("[%d] expected %#v\ngot\n%#v", i, row.err, err)
		}
	}

}

func TestFileExists(t *testing.T) {
	dir, _ := ioutil.TempDir("", "taskexists")
	createTestPkg(dir)
	table := []struct {
		dir    string
		task   string
		result bool
	}{
		{
			dir:    path.Join(dir, "script.d"),
			task:   "install",
			result: true,
		},
		{
			dir:    path.Join(dir, "5e7rszht.d"),
			task:   "install",
			result: false,
		},
	}

	for i, row := range table {
		result := fileExists(row.dir, row.task)
		if result != row.result {
			t.Fatalf("[%d] expected %t\ngot\n%t", i, row.result, result)
		}
	}
}

func createTestPkg(root string) {
	os.MkdirAll(path.Join(root, "script.d", "install"), os.ModePerm)
	os.MkdirAll(path.Join(root, "script.d", "errdir"), os.ModePerm)
	// task-dir: success
	ioutil.WriteFile(path.Join(root, "script.d", "install", "00-yolo"), []byte("#!/bin/bash\necho 'success1'; exit 0"), os.ModePerm)
	ioutil.WriteFile(path.Join(root, "script.d", "install", "02-aulu"), []byte("#!/bin/bash\necho 'success2'; exit 0"), os.ModePerm)
	// task-dir: error
	ioutil.WriteFile(path.Join(root, "script.d", "errdir", "02-aulu"), []byte("#!/bin/bash\necho 'errdir' >&2; exit 125"), os.ModePerm)
	// task-file: err
	ioutil.WriteFile(path.Join(root, "script.d", "error"), []byte("#!/bin/bash\necho 'err1'; exit 1"), os.ModePerm)
	// task-file: stderr
	ioutil.WriteFile(path.Join(root, "script.d", "stderror"), []byte("#!/bin/bash\necho 'err1231' >&2; exit 12"), os.ModePerm)
	// task-file: success
	ioutil.WriteFile(path.Join(root, "script.d", "success"), []byte("#!/bin/bash\necho 'succ2'; exit 0"), os.ModePerm)
}
