# Chill

`ionice` (from util-linux) ported to Go.

## Related projects

* [ion](https://github.com/xyproto/ion) is a fork of `ionice`, in 326 lines of C.
* [ionice](https://github.com/xyproto/ionice) is a Go module where the core functionality of the `ionice` utility has been ported to Go.

`chill` uses the [`ionice`](https://github.com/xyproto/ionice) Go module.

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
