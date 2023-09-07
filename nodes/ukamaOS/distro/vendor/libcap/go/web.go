// Progam web provides an example of a webserver using capabilities to
// bind to a privileged port, and then drop all capabilities before
// handling the first web request.
//
// This program cannot work reliably as a pure Go application without
// the equivalent of the Go runtime patch that adds a POSIX semantics
// wrapper around the system calls that change per-thread security
// state. A patch for the pure Go compiler/runtime to add this support
// is available here [2019-12-14]:
//
//    https://go-review.googlesource.com/c/go/+/210639/
//
// Until that patch, or something like it, is absorbed into the Go
// runtime the only way to get capabilities to work reliably on the Go
// runtime is to use something like libpsx via CGo to do capability
// setting syscalls in C with POSIX semantics. As of this build of the
// Go "kernel.org/pub/linux/libs/security/libcap/cap" package,
// courtesy of the "kernel.org/pub/linux/libs/security/libcap/psx"
// package, this is how things work.
//
// To set this up, compile and empower this binary as follows (read
// over the detail in the psx package description if this doesn't
// 'just' work):
//
//   go build web.go
//   sudo setcap cap_setpcap,cap_net_bind_service=p web
//   ./web --port=80
//
// Make requests using wget and observe the log of web:
//
//   wget -o/dev/null -O/dev/stdout localhost:80
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"syscall"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

var (
	port     = flag.Int("port", 0, "port to listen on")
	skipPriv = flag.Bool("skip", false, "skip raising the effective capability - will fail for low ports")
)

// ensureNotEUID aborts the program if it is running setuid something,
// or being invoked by root.  That is, the preparer isn't setting up
// the program correctly.
func ensureNotEUID() {
	euid := syscall.Geteuid()
	uid := syscall.Getuid()
	egid := syscall.Getegid()
	gid := syscall.Getgid()
	if uid != euid || gid != egid {
		log.Fatalf("go runtime is setuid uids:(%d vs %d), gids(%d vs %d)", uid, euid, gid, egid)
	}
	if uid == 0 {
		log.Fatalf("go runtime is running as root - cheating")
	}
}

// listen creates a listener by raising effective privilege only to
// bind to address and then lowering that effective privilege.
func listen(network, address string) (net.Listener, error) {
	if *skipPriv {
		return net.Listen(network, address)
	}

	orig := cap.GetProc()
	defer orig.SetProc() // restore original caps on exit.

	c, err := orig.Dup()
	if err != nil {
		return nil, fmt.Errorf("failed to dup caps: %v", err)
	}

	if on, _ := c.GetFlag(cap.Permitted, cap.NET_BIND_SERVICE); !on {
		return nil, fmt.Errorf("insufficient privilege to bind to low ports - want %q, have %q", cap.NET_BIND_SERVICE, c)
	}

	if err := c.SetFlag(cap.Effective, true, cap.NET_BIND_SERVICE); err != nil {
		return nil, fmt.Errorf("unable to set capability: %v", err)
	}

	if err := c.SetProc(); err != nil {
		return nil, fmt.Errorf("unable to raise capabilities %q: %v", c, err)
	}
	return net.Listen(network, address)
}

// Handler is used to abstract the ServeHTTP function.
type Handler struct{}

// ServeHTTP says hello from a single Go hardware thread and reveals
// its capabilities.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	runtime.LockOSThread()
	// Get some numbers consistent to the current execution, so
	// the returned web page demonstrates that the code execution
	// is bouncing around on different kernel thread ids.
	p := syscall.Getpid()
	t := syscall.Gettid()
	c := cap.GetProc()
	runtime.UnlockOSThread()

	log.Printf("Saying hello from proc: %d->%d, caps=%q", p, t, c)
	fmt.Fprintf(w, "Hello from proc: %d->%d, caps=%q\n", p, t, c)
}

func main() {
	flag.Parse()

	if *port == 0 {
		log.Fatal("please supply --port value")
	}

	ensureNotEUID()

	ls, err := listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("aborting: %v", err)
	}
	defer ls.Close()

	if !*skipPriv {
		if err := cap.ModeNoPriv.Set(); err != nil {
			log.Fatalf("unable to drop all privilege: %v", err)
		}
	}

	if err := http.Serve(ls, &Handler{}); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
