// chill is a port of ionice to Go

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/xyproto/gionice"
)

const versionString = "chill 1.2.0"

const usageString = `Usage:
 chill [options] -p <pid>...
 chill [options] -P <pgid>...
 chill [options] -u <uid>...
 chill [options] <command>

Set or get the I/O-scheduling class and priority of a process.

Options:
 -c, --class <class>    name or number of scheduling class,
                        0: none, 1: realtime, 2: best-effort, 3: idle
 -n, --classdata <num>  priority (0..7) in the specified scheduling class,
                        only for the realtime and best-effort classes
 -p, --pid <pid>...     act on these already running processes
 -P, --pgid <pgrp>...   act on already running processes in these groups
 -t, --ignore           ignore failures
 -N, --nice             set the niceness to 10
 -a, --adjustment <x>   adjust the nice priority with the given number
 -u, --uid <uid>...     act on already running processes owned by these users
 -s, --setnice <x>      set the process niceness

 -h, --help             display this help
 -V, --version          display version

For more details see chill(1).`

// Options is a struct containing information about all flags used by chill
type Options struct {
	Class      string `short:"c" long:"class" description:"name or number of scheduling class, 0: none, 1: realtime, 2: best-effort, 3: idle" choice:"0" choice:"1" choice:"2" choice:"3" choice:"none" choice:"realtime" choice:"best-effort" choice:"idle"`
	ClassData  int    `short:"n" long:"classdata" description:"priority (0..7) in the specified scheduling class, only for the realtime and best-effort classes" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5" choice:"6" choice:"7" choice:"8" choice:"9"`
	PID        int    `short:"p" long:"pid" description:"act on these already running processes" value-name:"PID"`
	PGID       int    `short:"P" long:"pgid" description:"act on already running processes in these groups" value-name:"PGID"`
	Ignore     bool   `short:"t" long:"ignore" description:"ignore failures"`
	Nice       bool   `short:"N" long:"nice" description:"also set niceness to 10"`
	UID        int    `short:"u" long:"uid" description:"act on already running processes owned by these users" value-name:"UID"`
	Help       bool   `short:"h" long:"help" description:"display this help"`
	Version    bool   `short:"V" long:"version" description:"display version"`
	Adjustment int    `short:"a" long:"adjustment" description:"niceness priority adjustment"`
	SetNice    int    `short:"s" long:"setnice" description:"set the process niceness"`
	Args       struct {
		Command []string
	}
}

func main() {
	opts := &Options{}
	parser := flags.NewParser(opts, flags.PassAfterNonOption|flags.PassDoubleDash)
	args, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var (
		hasClass      = parser.FindOptionByLongName("class").IsSet()
		hasClassData  = parser.FindOptionByLongName("classdata").IsSet()
		hasPID        = parser.FindOptionByLongName("pid").IsSet()
		hasPGID       = parser.FindOptionByLongName("pgid").IsSet()
		hasUID        = parser.FindOptionByLongName("uid").IsSet()
		hasAdjustment = parser.FindOptionByLongName("adjustment").IsSet()
		hasSetNice    = parser.FindOptionByLongName("setnice").IsSet()

		data            = 4
		set, which, who int
		ioclass         = gionice.IOPRIO_CLASS_BE
		tolerant        bool
	)

	if opts.Help {
		fmt.Println(usageString)
		os.Exit(0)
	}
	if opts.Version {
		fmt.Println(versionString)
		os.Exit(0)
	}
	if opts.Ignore {
		tolerant = true
	}
	if hasClassData {
		set |= 1
	}
	if hasClass {
		ioclass, err = gionice.Parse(opts.Class)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if ioclass < 0 {
			fmt.Fprintf(os.Stderr, "uknown scheduling class: '%s'\n", opts.Class)
		}
		set |= 2
	}
	if hasPID {
		if who != 0 {
			fmt.Fprintln(os.Stderr, "can handle only one of pid, pgid or uid at once")
			os.Exit(1)
		}
		which = opts.PID
		who = gionice.IOPRIO_WHO_PROCESS
	}
	if hasPGID {
		if who != 0 {
			fmt.Fprintln(os.Stderr, "can handle only one of pid, pgid or uid at once")
			os.Exit(1)
		}
		which = opts.PGID
		who = gionice.IOPRIO_WHO_PGRP
	}
	if hasUID {
		if who != 0 {
			fmt.Fprintln(os.Stderr, "can handle only one of pid, pgid or uid at once")
			os.Exit(1)
		}
		which = opts.UID
		who = gionice.IOPRIO_WHO_USER
	}

	// The functionality of "nice"
	if hasSetNice {
		gionice.SetNicePri(which, who, opts.SetNice)
	} else if hasAdjustment {
		currentPri, err := gionice.NicePri(which, who)
		if err != nil {
			// warning, can not get niceness priority
			fmt.Fprintf(os.Stderr, "can not get niceness: %v\n", err)
			os.Exit(1)
		}
		currentPri += opts.Adjustment
		gionice.SetNicePri(which, who, currentPri)
	} else if opts.Nice {
		gionice.SetNicePri(which, who, 10)
	}

	switch ioclass {
	case gionice.IOPRIO_CLASS_NONE:
		if (set&1) != 0 && !tolerant {
			// warning
			fmt.Fprintln(os.Stderr, "ignoring given cass data for none class")
		}
		data = 0
	case gionice.IOPRIO_CLASS_RT, gionice.IOPRIO_CLASS_BE:
		break
	case gionice.IOPRIO_CLASS_IDLE:
		if (set&1) != 0 && !tolerant {
			// just a warning, no exit
			fmt.Fprintln(os.Stderr, "ignoring given class data for idle class")
		}
		data = 7
	default:
		if !tolerant {
			// just a warning, no exit
			fmt.Fprintf(os.Stderr, "unknown prio class %d\n", ioclass)
		}
	}

	if set == 0 && which == 0 && len(args) == 0 {
		// chill without options, print the current ioprio
		gionice.Print(0, gionice.IOPRIO_WHO_PROCESS)
	} else if set == 0 && who != 0 {
		// chill -p|-P|-u ID [ID ...]
		gionice.Print(which, who)
		for _, id := range args {
			if n, err := strconv.Atoi(id); err == nil { // success, arg is a number
				which = n
				gionice.Print(which, who)
			}
		}
	} else if set != 0 && who != 0 {
		// chill -c CLASS -p|-P|-u ID [ID ...]
		if err := gionice.SetIDPri(which, ioclass, data, who); err != nil && !tolerant {
			fmt.Fprintln(os.Stderr, "ioprio_set failed", err)
			os.Exit(1)
		}
		for _, id := range args {
			if n, err := strconv.Atoi(id); err == nil { // success, arg is a number
				which = n

				if err := gionice.SetIDPri(which, ioclass, data, who); err != nil && !tolerant {
					fmt.Fprintln(os.Stderr, "ioprio_set failed", err)
					os.Exit(1)
				}
			}
		}
	} else if len(args) > 0 {
		// chill [-c CLASS] COMMAND
		if err := gionice.SetIDPri(0, ioclass, data, gionice.IOPRIO_WHO_PROCESS); err != nil && !tolerant {
			fmt.Fprintln(os.Stderr, "ioprio_set failed", err)
			os.Exit(1)
		}
		var argv0 string = args[0] // got to find the path first?
		argv0, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not find %s in PATH\n", args[0])
			os.Exit(1)
		}
		err = syscall.Exec(argv0, args, os.Environ())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute %s\n", argv0)
			os.Exit(1)
		}
		os.Exit(1)
	} else {
		fmt.Fprintln(os.Stderr, "bad usage\nTry 'chill --help' for more information.")
		os.Exit(1)
	}
}
