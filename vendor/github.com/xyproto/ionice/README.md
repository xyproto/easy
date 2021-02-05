# ionice

This is an extraction of the core parts of the `ionice` utility from `util-linux`, to a Go module.

The command line utility [chill](https://github.com/xyproto/chill), which is a drop-in replacement for the `ionice` utility, uses this module.

This package can be used by any Go program that wishes to run without hogging the I/O capabilities of the current system.

## General info

* Version: 1.0.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
