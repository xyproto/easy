package gionice

import (
	"syscall"
)

const (
	// From include/bits/resource.h
	PRIO_PROCESS = 0
	PRIO_PGRP    = 1
	PRIO_USER    = 2
)

// SetNicePri sets the IO priority for the given "which" (process, pgrp or user) and "who" (the ID),
// using the given io priority number.
func SetNicePri(which, who, ioprio int) error {
	return syscall.Setpriority(which, who, ioprio)
}

// NicePri returns the IO priority for the given "which" (process, pgrp or user) and "who" (the ID).
func NicePri(which, who int) (int, error) {
	return syscall.Getpriority(which, who)
}

// If permitted, set the given process group niceness to level 10.
// Pass in PGID 0 for the current process group.
func SetNice(pgid int) error {
	return SetNicePri(pgid, PRIO_PGRP, 10)
}

// If permitted, set the given process niceness to level 10.
// Pass in PID 0 for the current process.
func SetNicePID(pid int) error {
	return SetNicePri(pid, PRIO_PROCESS, 10)
}

// If permitted, set the given process group nicess to level -20.
// Pass in PGID 0 for the current process group.
func SetNaughty(pgid int) error {
	return SetNicePri(pgid, PRIO_PGRP, -20)
}

// If permitted, set the given process nicess to level -20.
// Pass in PID 0 for the current process.
func SetNaughtyPID(pid int) error {
	return SetNicePri(pid, PRIO_PROCESS, -20)
}

// Nice will try to set the current process group niceness to level 10.
func Nice() error {
	return SetNicePri(0, PRIO_PGRP, 10)
}

// Naughty will try to set the current process group niceness to level -20.
func Naughty() error {
	return SetNicePri(0, PRIO_PGRP, -20)
}
