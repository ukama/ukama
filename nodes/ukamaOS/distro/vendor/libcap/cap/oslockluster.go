// +build !go1.10

package cap

import "syscall"

// LaunchSupported indicates that is safe to return from a locked OS
// Thread and have that OS Thread be terminated by the runtime. The
// Launch functionality really needs to rely on the fact that an
// excess of runtime.LockOSThread() vs. runtime.UnlockOSThread() calls
// in a returning go routine will cause the underlying locked OSThread
// to terminate. That feature was added to the Go runtime in version
// 1.10.
//
// See these bugs for the discussion and feature assumed by the code
// in this Launch() functionality:
//
//   https://github.com/golang/go/issues/20395
//   https://github.com/golang/go/issues/20458
//
// A value of false for this constant causes cap.(*Launcher).Launch()
// to park the go routine used to perform the launch indefinitely so
// its kernel privilege state of the OS Thread locked to it does not
// pollute the rest of the runtime - yes, it leaks an OSThread. If
// this is a problem for your application you have two workarounds:
//
// 1) don't use cap.(*Launcher).Launch()
// 2) upgrade your Go toolchain to 1.10+
const LaunchSupported = false

// validatePA confirms that the pa.Sys entry is not incompatible with
// Launch.
func validatePA(pa *syscall.ProcAttr, chroot string) (bool, error) {
	return false, ErrNoLaunch
}
