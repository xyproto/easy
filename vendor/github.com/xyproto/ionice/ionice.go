// Package ionice contains code that has been ported from util-linux/ionice (GPL2 licensed)
package ionice

import (
	"fmt"
	"log"
	"strings"
	"syscall"
)

// PriClass represents an IO class, like "realtime" or "idle"
type PriClass int

const (
	IOPRIO_CLASS_NONE PriClass = 0
	IOPRIO_CLASS_RT   PriClass = 1
	IOPRIO_CLASS_BE   PriClass = 2
	IOPRIO_CLASS_IDLE PriClass = 3

	IOPRIO_WHO_PROCESS = 1
	IOPRIO_WHO_PGRP    = 2
	IOPRIO_WHO_USER    = 3

	IOPRIO_CLASS_SHIFT = 13
)

// SetPri sets the IO priority for the given "which" (process, pgrp or user) and "who" (the ID),
// using the given io priority number.
func SetPri(which, who int, ioprio uint) (uint, error) {
	r1, _, errNo := syscall.Syscall(syscall.SYS_IOPRIO_SET, uintptr(which), uintptr(who), uintptr(ioprio))
	if errNo != 0 {
		return uint(r1), errNo
	}
	return uint(r1), nil

}

// Pri returns the IO priority for the given "which" (process, pgrp or user) and "who" (the ID).
func Pri(which, who int) (uint, error) {
	r1, _, errNo := syscall.Syscall(syscall.SYS_IOPRIO_GET, uintptr(which), uintptr(who), uintptr(0))
	if errNo != 0 {
		return uint(r1), errNo
	}
	return uint(r1), nil
}

func priMask() uint {
	return (uint(1) << IOPRIO_CLASS_SHIFT) - 1
}

func priPriClass(mask uint) PriClass {
	return PriClass(mask >> IOPRIO_CLASS_SHIFT)
}

func priData(mask uint) uint {
	return mask & priMask()
}

func priValue(classn, data uint) uint {
	return ((classn << IOPRIO_CLASS_SHIFT) | data)
}

var to_prio = map[PriClass]string{
	IOPRIO_CLASS_NONE: "none",
	IOPRIO_CLASS_RT:   "realtime",
	IOPRIO_CLASS_BE:   "best-effort",
	IOPRIO_CLASS_IDLE: "idle",
}

// Parse converts a string containing either:
// "none", "realtime", best-effort" or "idle", to a corresponding IOPRIO_CLASS.
// will also handle "0", "1", "2" or "3"
// The parsing is case-insensitive, so "REALTIME" or "rEaLtImE" is also fine.
func Parse(ioprio string) (PriClass, error) {
	switch strings.ToLower(ioprio) {
	case "0", "none":
		return IOPRIO_CLASS_NONE, nil
	case "1", "realtime":
		return IOPRIO_CLASS_RT, nil
	case "2", "best-effort":
		return IOPRIO_CLASS_BE, nil
	case "3", "idle":
		return IOPRIO_CLASS_IDLE, nil
	}
	return 0, fmt.Errorf("could not parse %s as an IOPRIO_CLASS constant", ioprio)
}

// Print outputs the IO nice status for the given PID and "who"
func Print(pid, who int) {
	ioprio, err := Pri(who, pid)
	if err != nil {
		log.Fatalln("ioprio_get failed", err)
	}
	ioclass := priPriClass(ioprio)
	name := "unknown"
	to_prio_len := PriClass(len(to_prio))
	if ioclass >= 0 && ioclass < to_prio_len {
		name = to_prio[ioclass]
	}
	if ioclass != IOPRIO_CLASS_IDLE {
		fmt.Printf("%s: prio %d\n", name, priData(ioprio))
	} else {
		fmt.Printf("%s\n", name)
	}
}

func SetIDPri(which int, ioclass PriClass, data, who int, tolerant bool) {
	_, err := SetPri(who, which, priValue(uint(ioclass), uint(data)))
	if err != nil && !tolerant {
		log.Fatalln("ioprio_set failed", err)
	}
}
