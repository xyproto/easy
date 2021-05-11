# Chill

Chill implements the functionality of both `nice` and `ionice`.

It's also a drop-in replacement for `ionice` (from util-linux).

* Chill started out as a port of `ionice` to Go, but more functionality has been added since then.
* Chill can also be used to give applications increased I/O priority or niceness.
* Chill many be useful for running ie. Zoom or Chromium on desktop Linux, with a lower I/O priority.

## Differences from `ionice`

These flags are for adjusting the process niceness (from `nice` not `ionice`):

* `-N` or `--nice` can be used to **also** set the process niceness to 10 (same as `nice COMMAND`).
* `-s` or `--setnice` can be used to set the process niceness.
* `-a` or `--adjustment` can be used to adjust the process niceness by the given offset.

## Related projects

* `ionice` from util-linux.
* `nice` from coreutils.
* [ion](https://github.com/xyproto/ion) is a fork of `ionice`, in 326 lines of C.
* [gionice](https://github.com/xyproto/gionice) is a Go module where the core functionality of the `ionice` utility has been ported to Go.

## Why

This port exists mainly because I wanted to have a [Go module](https://github.com/xyproto/gionice) for changing the I/O priority of servers written in Go. It was relatively easy to add a port of the `ionice` utility as well, once that was done.

## Requirements

* A recent Go compiler.
* Linux.

## Build

    go build -mod=vendor

## Install

    install -Dm755 chill /usr/bin/chill
    gzip chill.1
    install -Dm644 chill.1.gz /usr/share/man/man1/chill.1.gz

## Build and install with the go command

    go get -u github.com/xyproto/chill

## General info

* Version: 1.2.0
* Licence: GPL2
