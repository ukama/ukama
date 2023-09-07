// Program posix is a test case to confirm that Go is capable of
// exhibiting posix semantics for system calls.
//
// This code is a template for two programs: posix.go and posix-cgo.go
// which are built by the Makefile to using sed.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
)

// main_here

func dumpStatus(testCase string, err error, filter, expect string) bool {
	fmt.Printf("%s [%v]:\n", testCase, err)
	var failed bool
	pid := syscall.Getpid()
	fs, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/task", pid))
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		tf := fmt.Sprintf("/proc/%s/status", f.Name())
		d, err := ioutil.ReadFile(tf)
		if err != nil {
			fmt.Println(tf, err)
			failed = true
			continue
		}
		lines := strings.Split(string(d), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, filter) {
				fails := line != expect
				failure := ""
				if fails {
					failed = fails
					failure = " (bad)"
				}
				fmt.Printf("%s %s%s\n", tf, line, failure)
				break
			}
		}
	}
	return failed
}

func ptest() {
	var err error
	var bad bool

	// egid setting
	bad = bad || dumpStatus("initial state", nil, "Gid:", "Gid:\t0\t0\t0\t0")
	err = syscall.Setegid(1001)
	bad = bad || dumpStatus("setegid(1001) state", err, "Gid:", "Gid:\t0\t1001\t0\t1001")
	err = syscall.Setegid(1002)
	bad = bad || dumpStatus("setegid(1002) state", err, "Gid:", "Gid:\t0\t1002\t0\t1002")
	err = syscall.Setegid(0)
	bad = bad || dumpStatus("setegid(0) state", err, "Gid:", "Gid:\t0\t0\t0\t0")

	// euid setting (no way back from this one)
	bad = bad || dumpStatus("initial euid", nil, "Uid:", "Uid:\t0\t0\t0\t0")
	err = syscall.Seteuid(1)
	bad = bad || dumpStatus("seteuid(1)", err, "Uid:", "Uid:\t0\t1\t0\t1")
	err = syscall.Seteuid(0)
	bad = bad || dumpStatus("seteuid(0)", err, "Uid:", "Uid:\t0\t0\t0\t0")

	// gid setting
	bad = bad || dumpStatus("initial state", nil, "Gid:", "Gid:\t0\t0\t0\t0")
	err = syscall.Setgid(1001)
	bad = bad || dumpStatus("setgid(1001) state", err, "Gid:", "Gid:\t1001\t1001\t1001\t1001")
	err = syscall.Setgid(1002)
	bad = bad || dumpStatus("setgid(1002) state", err, "Gid:", "Gid:\t1002\t1002\t1002\t1002")
	err = syscall.Setgid(0)
	bad = bad || dumpStatus("setgid(0) state", err, "Gid:", "Gid:\t0\t0\t0\t0")

	// groups setting
	bad = bad || dumpStatus("initial groups", nil, "Groups:", "Groups:\t0 ")
	err = syscall.Setgroups([]int{0, 1, 2, 3})
	bad = bad || dumpStatus("setgroups(0,1,2,3)", err, "Groups:", "Groups:\t0 1 2 3 ")
	err = syscall.Setgroups([]int{3, 2, 1})
	bad = bad || dumpStatus("setgroups(2,3,1)", err, "Groups:", "Groups:\t1 2 3 ")
	err = syscall.Setgroups(nil)
	bad = bad || dumpStatus("setgroups(nil)", err, "Groups:", "Groups:\t ")
	err = syscall.Setgroups([]int{0})
	bad = bad || dumpStatus("setgroups(0)", err, "Groups:", "Groups:\t0 ")

	// regid setting
	bad = bad || dumpStatus("initial state", nil, "Gid:", "Gid:\t0\t0\t0\t0")
	err = syscall.Setregid(1001, 0)
	bad = bad || dumpStatus("setregid(1001) state", err, "Gid:", "Gid:\t1001\t0\t0\t0")
	err = syscall.Setregid(0, 1002)
	bad = bad || dumpStatus("setregid(1002) state", err, "Gid:", "Gid:\t0\t1002\t1002\t1002")
	err = syscall.Setregid(0, 0)
	bad = bad || dumpStatus("setregid(0) state", err, "Gid:", "Gid:\t0\t0\t0\t0")

	// reuid setting
	bad = bad || dumpStatus("initial state", nil, "Uid:", "Uid:\t0\t0\t0\t0")
	err = syscall.Setreuid(1, 0)
	bad = bad || dumpStatus("setreuid(1,0) state", err, "Uid:", "Uid:\t1\t0\t0\t0")
	err = syscall.Setreuid(0, 2)
	bad = bad || dumpStatus("setreuid(0,2) state", err, "Uid:", "Uid:\t0\t2\t2\t2")
	err = syscall.Setreuid(0, 0)
	bad = bad || dumpStatus("setreuid(0) state", err, "Uid:", "Uid:\t0\t0\t0\t0")

	// resgid setting
	bad = bad || dumpStatus("initial state", nil, "Gid:", "Gid:\t0\t0\t0\t0")
	err = syscall.Setresgid(1, 0, 2)
	bad = bad || dumpStatus("setresgid(1,0,2) state", err, "Gid:", "Gid:\t1\t0\t2\t0")
	err = syscall.Setresgid(0, 2, 1)
	bad = bad || dumpStatus("setresgid(0,2,1) state", err, "Gid:", "Gid:\t0\t2\t1\t2")
	err = syscall.Setresgid(0, 0, 0)
	bad = bad || dumpStatus("setresgid(0) state", err, "Gid:", "Gid:\t0\t0\t0\t0")

	// resuid setting
	bad = bad || dumpStatus("initial state", nil, "Uid:", "Uid:\t0\t0\t0\t0")
	err = syscall.Setresuid(1, 0, 2)
	bad = bad || dumpStatus("setresuid(1,0,2) state", err, "Uid:", "Uid:\t1\t0\t2\t0")
	err = syscall.Setresuid(0, 2, 1)
	bad = bad || dumpStatus("setresuid(0,2,1) state", err, "Uid:", "Uid:\t0\t2\t1\t2")
	err = syscall.Setresuid(0, 0, 0)
	bad = bad || dumpStatus("setresuid(0) state", err, "Uid:", "Uid:\t0\t0\t0\t0")

	// uid setting (no way back from this one)
	bad = bad || dumpStatus("initial uid", nil, "Uid:", "Uid:\t0\t0\t0\t0")
	err = syscall.Setuid(1)
	bad = bad || dumpStatus("setuid(1)", err, "Uid:", "Uid:\t1\t1\t1\t1")
	err = syscall.Setuid(0)
	bad = bad || dumpStatus("setuid(0)", err, "Uid:", "Uid:\t1\t1\t1\t1")

	if bad {
		log.Print("TEST FAILED")
		os.Exit(1)
	}
	log.Print("TEST PASSED")
}
