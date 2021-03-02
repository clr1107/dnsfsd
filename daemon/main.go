package main

import (
	"fmt"
	"github.com/clr1107/dnsfsd/pkg/data/cache"
	"github.com/clr1107/dnsfsd/pkg/data/config"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/clr1107/dnsfsd/daemon/logger"
	"github.com/clr1107/dnsfsd/daemon/server"
	"github.com/clr1107/dnsfsd/pkg/rules"
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

	if err := config.InitConfig(); err != nil {
		fmt.Printf("main() init config: %v\n", err)
		return
	}

	logPath := viper.Sub("log").GetString("path")
	port := viper.Sub("server").GetInt("port")
	forwards := viper.Sub("dns").GetStringSlice("forwards")
	verbose := viper.Sub("log").GetBool("verbose")
	cacheTTL := config.GetCacheTime()

	if err := log.Init(logPath); err != nil {
		fmt.Printf("main() init loggers: %v\n", err)
		os.Exit(1)
	}

	log.Log(strings.Repeat("=", 80))

	loadedRules, err := loadRules()
	if err != nil {
		log.LogFatal("main() loading rules: %v", err)
	} else {
		log.Log("loaded %v rules", loadedRules.Size())
	}

	dnsCache, err := cache.DNSCacheFromFile(cacheTTL, "/etc/dnsfsd/dns.cache")
	if err != nil {
		log.LogErr("could not load dns cache file, creating new DNSCache")
		dnsCache = cache.NewDNSCache(cacheTTL)
	} else {
		log.Log("loaded %v requests from the disk cache", dnsCache.Size())
	}

	srv := server.NewServer(port, server.NewHandler(loadedRules, dnsCache, forwards, verbose, log))
	spawnSignalRoutine(srv)

	go func() {
		for err := range srv.Handler.ErrorChannel {
			log.LogErr("server error listener: %v", err)
		}
	}()

	log.Log("starting listening on port %v with %v servers (verbose: %v)", port, len(forwards), verbose)
	if err := srv.Server.ListenAndServe(); err != nil {
		log.LogFatal("main() starting server: %v", err)
	}
}
