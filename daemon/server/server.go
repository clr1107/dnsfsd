package server

import (
	"regexp"
	"strconv"

	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/miekg/dns"
)

type DNSFSServer struct {
	Port    int
	server  *dns.Server
	handler *handler
}

func (s *DNSFSServer) Start() error {
	s.server = &dns.Server{Addr: ":" + strconv.Itoa(s.Port), Net: "udp"}

	files, err := persistence.LoadAllPatternFiles("/etc/dnsfsd/patterns")
	if err != nil {
		return err
	}

	patterns := persistence.CollectAllPatterns(files)
	s.handler = newHandler(patterns) // just copy it accross to our own struct.

	s.server.Handler = s.handler
	return s.server.ListenAndServe()
}

type handler struct {
	patterns []*regexp.Regexp
	cache    *persistence.SimpleCache
}

func newHandler(patterns []*regexp.Regexp) *handler {
	return &handler{patterns, persistence.NewSimpleCache(120000)}
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	// todo
}
