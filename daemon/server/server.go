package server

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/miekg/dns"
)

type DNSFSServer struct {
	Port    int
	Server  *dns.Server
	Handler *DNSFSHandler
}

func (s *DNSFSServer) Init() error {
	s.Server = &dns.Server{Addr: ":" + strconv.Itoa(s.Port), Net: "udp"}

	files, err := persistence.LoadAllPatternFiles("/etc/dnsfsd/patterns")
	if err != nil {
		return err
	}

	patterns := persistence.CollectAllPatterns(files)
	s.Handler = newHandler(patterns) // just copy it accross to our own struct.

	s.Server.Handler = s.Handler
	return nil
}

func (s *DNSFSServer) Shutdown() error {
	s.Handler.cache.Clear()
	close(s.Handler.ErrorChannel)

	return s.Server.Shutdown()
}

type DNSFSHandler struct {
	patterns     []*regexp.Regexp
	cache        *persistence.SimpleCache
	ErrorChannel chan error
}

func newHandler(patterns []*regexp.Regexp) *DNSFSHandler {
	return &DNSFSHandler{patterns, persistence.NewSimpleCache(-1), make(chan error)}
}

// true => sink; false => nothing found
func (h *DNSFSHandler) checkPatterns(domain string) bool {
	for _, pattern := range h.patterns {
		if pattern.MatchString(domain) {
			return true
		}
	}

	return false
}

// returns whether to sink or not based on cache and pattern matching
func (h *DNSFSHandler) check(domain string) bool {
	if h.cache.Contains(domain) {
		if val, ok := h.cache.Get(domain).(bool); ok {
			return val
		}

		h.cache.Remove(domain) // for some reason not a bool?
	}

	if h.checkPatterns(domain) {
		h.cache.PutDefault(domain, true)
		return true
	}

	return false
}

func (h *DNSFSHandler) forward(w *dns.ResponseWriter, r *dns.Msg, dnsAddress string) {
	c := new(dns.Client)
	x, _, err := c.Exchange(r, dnsAddress)

	if strings.Contains(r.Question[0].Name, "errorpls") {
		err = errors.New("random err")
	}

	if err != nil || x == nil {
		if x == nil {
			err = fmt.Errorf("after forwarding to '%v' the message response was nil", dnsAddress)
		}

		h.ErrorChannel <- err

		x = &dns.Msg{}
		x.SetReply(r)
		x.Authoritative = r.Authoritative
	}

	(*w).WriteMsg(x)
}

func (h *DNSFSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	aType := r.Question[0].Qtype == dns.TypeA

	if aType {
		msg := dns.Msg{}
		msg.SetReply(r)
		msg.Authoritative = true

		domain := msg.Question[0].Name
		if h.check(domain) {
			w.WriteMsg(&msg) // just sink right now
			return
		}
	}

	go h.forward(&w, r, "1.1.1.1:53")
}
