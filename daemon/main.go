package main

import (
	"fmt"

	"github.com/clr1107/dnsfsd/daemon/server"
)

func main() {
	srv := &server.DNSFSServer{Port: 53}

	if err := srv.Start(); err != nil {
		fmt.Printf("starting error: %v\n", err)
		return
	}

	println("serving...")
}
