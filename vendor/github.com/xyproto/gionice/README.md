# gionice

This is an port of the core parts of the `ionice` utility from `util-linux`, to a Go module.

The command line utility [chill](https://github.com/xyproto/chill), which is a drop-in replacement for the `ionice` utility, uses this module.

This package can be used by any Go program that wishes to run without hogging the I/O capabilities of the current system.

## Example use

To make your own Go program run as "idle" and not hog the I/O capabilities of the system, simply call `ionice.SetIdle(0)`:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xyproto/gionice"
)

func main() {
	// Make the current process "idle" (level 7)
	gionice.SetIdle(0)

	// Write to a file and delete then delete it, repeatedly
	for {
		fmt.Println("TICK")
		_ = ioutil.WriteFile("frenetic.dat", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0644)
		fmt.Println("TOCK")
		_ = os.Remove("frenetic.dat")
	}
}
```

By using `iotop` it's easy to check that the process PRIO is now `idle`.

## General info

* Version: 1.2.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
