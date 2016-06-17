# harry

> Go ahead, make my day.

Automatically runs `make` when its prerequisites change (or, become "dirty").

## Usage

Pretend you're running `make`, but type `harry` instead. It will do an initial
build if necessary, and continuously watch for changes in the files and
directories described by the specified target (plus the `Makefile` itself).

All arguments are forwarded along to `make`, so if you want to continuously run
a particular target, just do `harry TARGET`. Same goes for forwarding
parameters (`harry FOO=1`).

## Does it work?

> Do you feel lucky, punk?

Filesystem-watching interfaces are notoriously inconsistent across platforms.
Harry should at least compile on your platform, but it may not work very well.
For example, on Darwin it has to resort to watching a ton of files, and you may
hit a resource limit. If this happens you can increase the limit with `ulimit
-n 65535`.

Linux should work fine. No idea about Windows.

## Dependencies

Just `make`. Runs `make -dnr` to figure out your prerequisites, and uses
`fsnotify` to watch for changes. This makes `harry` super portable.
