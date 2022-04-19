package psx

import (
	"syscall"
	"testing"
)

func TestSyscall3(t *testing.T) {
	want := syscall.Getpid()
	if got, _, err := Syscall3(syscall.SYS_GETPID, 0, 0, 0); err != 0 {
		t.Errorf("failed to get PID via libpsx: %v", err)
	} else if int(got) != want {
		t.Errorf("pid mismatch: got=%d want=%d", got, want)
	}
	if got, _, err := Syscall3(syscall.SYS_CAPGET, 0, 0, 0); err != 14 {
		t.Errorf("malformed capget returned %d: %v (want 14: %v)", err, err, syscall.Errno(14))
	} else if ^got != 0 {
		t.Errorf("malformed capget did not return -1, got=%d", got)
	}
}

func TestSyscall6(t *testing.T) {
	want := syscall.Getpid()
	if got, _, err := Syscall6(syscall.SYS_GETPID, 0, 0, 0, 0, 0, 0); err != 0 {
		t.Errorf("failed to get PID via libpsx: %v", err)
	} else if int(got) != want {
		t.Errorf("pid mismatch: got=%d want=%d", got, want)
	}
	if got, _, err := Syscall6(syscall.SYS_CAPGET, 0, 0, 0, 0, 0, 0); err != 14 {
		t.Errorf("malformed capget errno %d: %v (want 14: %v)", err, err, syscall.Errno(14))
	} else if ^got != 0 {
		t.Errorf("malformed capget did not return -1, got=%d", got)
	}
}
