# gionice

This is an port of the core parts of the `ionice` utility from `util-linux`, to a Go module, without using `cgo`.

The command line utility [chill](https://github.com/xyproto/chill) (a drop-in replacement for `ionice`), uses this module.

This package can be used by any Go program that wishes to run without taking up the I/O capabilities of the current system.

## Example use

To make your own Go program run as "idle" and not hog the I/O capabilities of the system, simply call `ionice.Idle()`:

```go
package main

import (
	"io/ioutil"
	"os"

	"github.com/xyproto/gionice"
)

func main() {
	// Make the current process group priority to be "idle" (level 7)
	gionice.Idle()

	// Generate I/O activity
	for {
		_ = ioutil.WriteFile("frenetic.dat", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0644)
		_ = os.Remove("frenetic.dat")
	}
}
```

By using `iotop` it's easy to check that the process PRIO is now `idle`.

## General info

* Version: 1.3.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
