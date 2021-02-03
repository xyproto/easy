# ionice

This is a port and rewrite of ionice (from util-linux, GPL2 licensed) to a Go module and a `gionice` utility.

The development version of `gionice` can be installed with:

    go get -u github.com/xyproto/ionice/cmd/gionice

The `ionice` Go package can be used by Go programs that wishes to run without hogging the IO capabilities of the current system.

Even though this code is based on solid and well tested C code, the current implementation needs more testing.

## General info

* Version: 0.9.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
