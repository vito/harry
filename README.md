# harry

> Go ahead, make my day.

Automatically runs `make` when its prerequisites change.

## Usage

Pretend you're running `make`, but type `harry` instead. It will do an initial
build if necessary, and continuously watch for changes in the files and
directories described by the specified target (plus the `Makefile` itself).

All arguments are forwarded along to `make`, so if you want to continuously run
a particular target, just do `harry TARGET`. Same goes for forwarding
parameters (`harry FOO=1`).

## Dependencies

Just `make`. Runs `make -dnr` to figure out your prerequisites, and uses
`fsnotify` to watch for changes. This makes `harry` super portable.
