package main

import (
	"flag"
	"fmt"
	cfg "lwm2m-gateway/pkg/config"
	"lwm2m-gateway/pkg/iface"
	"lwm2m-gateway/pkg/lwm2m"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Usage
func usageError() {
	fmt.Println("Usage: lwm2m-gateway -cfgPath <Path> ")
	fmt.Println(" cfgPath is only one optional argument accepted.")
	fmt.Println(" If not provided /etc/config is a default path used to search lwm2m-gateway.json ")
}

var serviceName = "LwM2MGateway"

// main
func main() {

	// Log level
	log.SetLevel(log.TraceLevel)
	log.Infoln("LwM2MGateway:: Starting services " + serviceName + "...!!\n")

	cfgPath := "/etc/config/lwm2m-gateway/"
	// Process Arguments
	if len(os.Args) > 1 {
		cfgPath = processArgs(os.Args)
	}

	// LoadConfig
	err := cfg.LoadConfig("lwm2m-gateway", "json", cfgPath)
	if err != nil {
		log.Errorf("LwM2MGateway:: LwM2MGateway:: Failed to load config. Err: %s", err.Error())
		os.Exit(1)
	}

	// Reset log level based on config
	log.SetLevel(getLogLevel())

	// Makes sure connection is closed when service exits.
	handleSigterm(func() {
		iface.Stop()
	})

	// Parse model
	//lwm2m.CheckIfOperationAvailOnResourcesId("3328", 5700, lwm2m.WRITE)

	// Initialize msgbus
	iface.Start()

	// Start receiver for lwm2m
	lwm2m.Receiver()

	for {
		time.Sleep(10 * time.Second)
	}

}

// Process commnd line argument
func processArgs(args []string) string {
	var cfgPathPtr *string
	if strings.HasPrefix(args[1], "-cfgPath") {
		cfgPathPtr = flag.String("cfgPath", "/etc/config", "Path to search config file lwm2m-gateway.json")
		flag.Parse()
	} else {
		usageError()
		os.Exit(1)
	}
	return *cfgPathPtr
}

// Get logger level
func getLogLevel() log.Level {
	var level log.Level

	switch cfg.Config.LogLevel {
	case "trace", "TRACE":
		level = log.TraceLevel
	case "debug", "DEBUG":
		level = log.DebugLevel
	case "info", "INFO":
		level = log.InfoLevel
	case "panic", "PANIC":
		level = log.TraceLevel
	case "fatal", "FATAL":
		level = log.FatalLevel
	case "error", "Error":
		level = log.ErrorLevel
	default:
		level = log.DebugLevel
	}

	return level
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()

}
