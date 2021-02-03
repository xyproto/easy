// Package ionice contains code that has been ported from util-linux/ionice (GPL2 licensed)
package ionice

import (
	"fmt"
	"log"
	"strings"
	"syscall"
)

const (
	IOPRIO_CLASS_NONE = 0
	IOPRIO_CLASS_RT   = 1
	IOPRIO_CLASS_BE   = 2
	IOPRIO_CLASS_IDLE = 3

	IOPRIO_WHO_PROCESS = 1
	IOPRIO_WHO_PGRP    = 2
	IOPRIO_WHO_USER    = 3

	IOPRIO_CLASS_SHIFT = 13
)

type IOPRIO_CLASS int

// SetPri sets the IO priority for the given which (process, pgrp or user) and who (the ID),
// using the given io priority number.
func SetPri(which, who int, ioprio uint) (uint, error) {
	r1, _, errNo := syscall.Syscall(syscall.SYS_IOPRIO_SET, uintptr(which), uintptr(who), uintptr(ioprio))
	var err error
	if errNo != 0 {
		err = errNo
	}
	return uint(r1), err
}

// Pri returns the IO priority for the given which (process, pgrp or user) and who (the ID).
func Pri(which, who int) (uint, error) {
	r1, _, errNo := syscall.Syscall(syscall.SYS_IOPRIO_GET, uintptr(which), uintptr(who), uintptr(0))
	var err error
	if errNo != 0 {
		err = errNo
	}
	// TODO: r1 or r2?
	return uint(r1), err
}

func IOPRIO_PRIO_MASK() uint {
	return (uint(1) << IOPRIO_CLASS_SHIFT) - 1
}

func IOPRIO_PRIO_CLASS(mask uint) IOPRIO_CLASS {
	return IOPRIO_CLASS(mask >> IOPRIO_CLASS_SHIFT)
}

func IOPRIO_PRIO_DATA(mask uint) uint {
	return mask & IOPRIO_PRIO_MASK()
}

func IOPRIO_PRIO_VALUE(classn, data uint) uint {
	return ((classn << IOPRIO_CLASS_SHIFT) | data)
}

var to_prio = map[IOPRIO_CLASS]string{
	IOPRIO_CLASS_NONE: "none",
	IOPRIO_CLASS_RT:   "realtime",
	IOPRIO_CLASS_BE:   "best-effort",
	IOPRIO_CLASS_IDLE: "idle",
}

// Parse converts a case-insensitive string containing either:
// "none", "realtime", best-effort" or "idle", to a corresponding IOPRIO_CLASS.
func Parse(str string) (IOPRIO_CLASS, error) {
	for k, v := range to_prio {
		if strings.ToLower(str) == strings.ToLower(v) {
			return k, nil
		}
	}
	return 0, fmt.Errorf("could not parse %s to an IOPRIO_CLASS constant", str)
}

func Print(pid, who int) {
	ioprio, err := Pri(who, pid)
	if err != nil {
		log.Fatalln("ioprio_get failed", err)
	}
	ioclass := IOPRIO_PRIO_CLASS(ioprio)
	name := "unknown"
	last_index := IOPRIO_CLASS(len(to_prio) - 1)
	if ioclass >= 0 && ioclass <= last_index {
		name = to_prio[ioclass]
	}
	if ioclass != IOPRIO_CLASS_IDLE {
		fmt.Printf("%s: prio %d\n", name, IOPRIO_PRIO_DATA(ioprio))
	} else {
		fmt.Printf("%s\n", name)
	}
}

func SetIDPri(which int, ioclass IOPRIO_CLASS, data, who int, tolerant bool) {
	_, err := SetPri(who, which, IOPRIO_PRIO_VALUE(uint(ioclass), uint(data)))
	if err != nil && !tolerant {
		log.Fatalln("ioprio_set failed", err)
	}
}
