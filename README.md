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
│   ├── build
│   │   └── 00-build
│   ├── install
│   │   ├── 00-buildtools
│   │   └── 01-deps
│   └── releas
```

You might have noticed, that `install` is a directory and `release`  is a file. If you run `qwerty run install` then all scripts inside the that folder are executed in-order. The actual exec command is `bash <file>` while bash has to be in yourt $PATH. all other env vars should be available.

#### Help

```
$ qwerty --help
NAME:
   qwerty - A new cli application

USAGE:
   qwerty [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     run, r   looks for a script.d directory and runs the specified command provided by args.
     init, i  scaffold a script.d directory
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value    set the working directory
   --help, -h     show help
   --version, -v  print the version
   ```