package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/clr1107/dnsfsd/daemon/logger"
	"github.com/clr1107/dnsfsd/daemon/server"
	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/spf13/viper"
)

var (
	log *logger.Logger = &logger.Logger{}
)

func loadPatterns() ([]*regexp.Regexp, error) {
	files, err := persistence.LoadAllPatternFiles("/etc/dnsfsd/patterns")

	if err != nil {
		return nil, err
	}

	return persistence.CollectAllPatterns(files), nil
}

func spawnSignalRoutine(srv *server.DNSFSServer) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		log.Log("shutting down...")

		if err := srv.Shutdown(); err != nil {
			log.LogFatal("signal listener shutting down: %v", err)
		}

		os.Exit(0)
	}()
}

func main() {
	println(strings.Repeat("=", 80))

	if err := persistence.InitConfig(); err != nil {
		fmt.Printf("main() init config: %v\n", err)
		return
	}

	logPath := viper.GetString("log")
	port := viper.GetInt("port")
	forwards := viper.GetStringSlice("forwards")
	verbose := viper.GetBool("verbose")

	if err := log.Init(logPath); err != nil {
		fmt.Printf("main() init loggers: %v\n", err)
		os.Exit(1)
	}

	patterns, err := loadPatterns()
	if err != nil {
		log.LogFatal("main() loading patterns: %v", err)
	} else {
		log.Log("loaded %v patterns", len(patterns))
	}

	srv := server.NewServer(port, server.NewHandler(patterns, forwards, verbose, log))

	spawnSignalRoutine(srv)

	go func() {
		for err := range srv.Handler.ErrorChannel {
			log.LogErr("server error listener: %v", err)
		}
	}()

	log.Log("starting to listen on port %v (verbose: %v)", port, verbose)
	if err := srv.Server.ListenAndServe(); err != nil {
		log.LogFatal("main() starting server: %v", err)
	}
}
