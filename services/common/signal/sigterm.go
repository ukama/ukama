package sig

import (
	"os"
	"os/signal"
	"syscall"
)

/* Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting. */
func HandleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()

}
