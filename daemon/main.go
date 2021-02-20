package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/clr1107/dnsfsd/daemon/server"
)

func main() {
	// this will obv be touched up a little.
	srv := &server.DNSFSServer{Port: 53}

	if err := srv.Init(); err != nil {
		fmt.Printf("init error: %v\n", err)
		return
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, os.Kill)

	go func() {
		<-signalChannel
		fmt.Printf("\nShutting down...\n")

		if err := srv.Shutdown(); err != nil {
			fmt.Printf("Error whilst shutting down: %v\n", err)
		}

		os.Exit(0)
	}()

	go func() {
		for error := range srv.Handler.ErrorChannel {
			fmt.Printf("Server Error: %v\n", error)
		}
	}()

	if err := srv.Server.ListenAndServe(); err != nil {
		fmt.Printf("server error whilst starting: %v\n", err)
		return
	}
}
