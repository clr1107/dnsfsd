package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/clr1107/dnsfsd/daemon/logger"
	"github.com/clr1107/dnsfsd/daemon/server"
	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/clr1107/dnsfsd/pkg/persistence/rules"
	"github.com/spf13/viper"
)

var (
	log *logger.Logger = &logger.Logger{}
)

func loadRules() (*rules.RuleSet, error) {
	files, err := rules.LoadAllRuleFiles("/etc/dnsfsd/rules")

	if err != nil {
		return nil, err
	}

	return rules.CollectAllRules(files), nil
}

func spawnSignalRoutine(srv *server.DNSFSServer) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		log.Log("interrupt signal; shutting down...")

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

	rules, err := loadRules()
	if err != nil {
		log.LogFatal("main() loading rules: %v", err)
	} else {
		log.Log("loaded %v rules", rules.Size())
	}

	srv := server.NewServer(port, server.NewHandler(rules, forwards, verbose, log))

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
