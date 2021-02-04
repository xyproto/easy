# ionice

This is a port and rewrite of `ionice` from `util-linux` to a Go module.

There's a command line utility named `chill` available for installation with:

    go get -u github.com/xyproto/chill

This will download and install the development version of [Chill](https://github.com/xyproto/chill), which uses this Go module.

The `ionice`  package can be used by any Go program that wishes to run without hogging the IO capabilities of the current system.

## General info

* Version: 1.0.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
