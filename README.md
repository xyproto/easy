# Chill

Chill is the `ionice` utility (from util-linux), ported to Go. It's a drop-in replacement.

It can be used for not letting applications use too much of the disk or network capacity.

This many be useful for running ie. Zoom or Chromium on desktop Linux.

It can also be used to give applications increased I/O priority.

Also, being able to launch programs like `nice chill chromium` is just more fun.

## Differences from `ionice`

* `-a` can be used to also set the process niceness to 10 (same as `nice COMMAND`).

## Related projects

* [ion](https://github.com/xyproto/ion) is a fork of `ionice`, in 326 lines of C.
* [gionice](https://github.com/xyproto/gionice) is a Go module where the core functionality of the `ionice` utility has been ported to Go.

`chill` uses the [`gionice`](https://github.com/xyproto/gionice) Go module.

## Why

This port exists mainly because I wanted to have a Go module for changing the I/O priority of servers written in Go. It was relatively easy to add a port of the `ionice` utility as well, once that was done.

The Go executable is slightly larger than one produced in C, but might provide additional memory safety.

## Requirements

Just Go and Linux.

## Build

    go build -mod=vendor

## Install

    install -Dm755 chill /usr/bin/chill
    gzip chill.1
    install -Dm644 chill.1.gz /usr/share/man/man1/chill.1.gz

## Build and install with the go command

    go get -u github.com/xyproto/chill

## General info

* Version: 1.0.0
* Licence: GPL2
