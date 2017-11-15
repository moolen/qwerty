### QWERTY

#### Usage

Create a script.d directory in your repository. E.g. like this: 
```
.
├── build.bash
├── glide.yaml
├── main.go
├── README.md
├── script.d
│   ├── build     <-- a task
│   │   └── 00-build
│   ├── install   <-- a task
│   │   ├── 00-buildtools
│   │   └── 01-deps
│   └── release   <-- a task
```

You might have noticed, that `install` is a directory and `release`  is a file. If you run `qwerty run install` then all scripts inside the that folder are executed in-order. The actual exec command is `bash <file>` while bash has to be in yourt $PATH. all other env vars should be available.

To run a task issue `qwerty run install` in your terminal. That will call `00-buildtools` and then `01-deps`.

You can also run a `mux` command that executes tasks within subdirectories subsequently. E.g. issue `qwerty mux prepare build run` that will first `prepare`, then `build` and finally `run` them all.

```
.
├── project1
│   └── script.d
│       ├── build
│       │   ├── 00-lint
│       │   ├── 01-artifact1
│       ├── prepare
.       │   ├── 00-deps
.       │   └── 01-git-hooks
.       └── run
.           └── 00-spawn
└── project12
    └── script.d
        ├── build
        │   ├── 00-lint
        │   ├── 01-artifact1
        ├── prepare
        │   ├── 00-deps
        │   └── 01-git-hooks
        └── run
            └── 00-spawn
```

#### Help

```
$ qwerty --help
NAME:
   qwerty - A new cli application

USAGE:
   qwerty [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     run, r   runs the specified command, you may specify multiple commands.
     mux, m   multiplex commands in subdirectories
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value    set the working directory
   --debug        enable debug
   --help, -h     show help
   --version, -v  print the version
   ```