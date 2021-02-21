package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/clr1107/dnsfsd/daemon/logger"
	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/miekg/dns"
)

type DNSFSServer struct {
	Port    int
	Server  *dns.Server
	Handler *DNSFSHandler
}

func NewServer(port int, handler *DNSFSHandler) *DNSFSServer {
	s := &DNSFSServer{Port: port}

	s.Server = &dns.Server{Addr: ":" + strconv.Itoa(s.Port), Net: "udp"}
	s.Server.Handler = handler
	s.Handler = handler

	return s
}

func (s *DNSFSServer) Shutdown() error {
	s.Handler.cache.Clear()
	close(s.Handler.ErrorChannel)

	return s.Server.Shutdown()
}

type DNSFSHandler struct {
	rules        *[]persistence.IRule
	cache        *persistence.SimpleCache
	forwards     []string
	ErrorChannel chan error
	verbose      bool
	logger       *logger.Logger
}

func NewHandler(rules *[]persistence.IRule, forwards []string, verbose bool, logger *logger.Logger) *DNSFSHandler {
	return &DNSFSHandler{rules, persistence.NewSimpleCache(-1), forwards, make(chan error), verbose, logger}
}

// true => sink; false => nothing found
func (h *DNSFSHandler) checkRules(domain string) bool {
	for _, rule := range *h.rules {
		if rule.Match(domain) {
			if h.verbose {
				h.logger.Log("rule: '%v'", rule)
			}
			return true
		}
	}

	return false
}

// returns whether to sink or not based on cache and rule matching
func (h *DNSFSHandler) check(domain string) bool {
	if h.cache.Contains(domain) {
		if val, ok := h.cache.Get(domain).(bool); ok {
			if h.verbose && val {
				h.logger.Log("(%v) in cache => sinkhole", domain)
			}

			return val
		}

		if h.verbose {
			h.logger.LogErr("(%v) in cache, non-boolean value, removing", domain)
		}

		h.cache.Remove(domain) // for some reason not a bool?
	}

	if h.checkRules(domain) {
		if h.verbose {
			h.logger.Log("(%v) matched rule(s), putting in cache => sink", domain)
		}

		h.cache.PutDefault(domain, true)
		return true
	} else if h.verbose {
		h.logger.Log("(%v) matched no rules, putting in cache => pass", domain)
	}

	h.cache.PutDefault(domain, false)
	return false
}

func (h *DNSFSHandler) forward(w *dns.ResponseWriter, r *dns.Msg, dnsAddress string) error {
	if h.verbose {
		h.logger.Log("(%v) forwarding to %v", r.Question[0].Name, strings.Replace(dnsAddress, ":", "#", 1))
	}

	c := new(dns.Client)
	x, _, err := c.Exchange(r, dnsAddress)

	if err != nil || x == nil {
		if x == nil {
			err = fmt.Errorf("after forwarding to '%v' the message response was nil", dnsAddress)
		}

		h.ErrorChannel <- err

		x = &dns.Msg{}
		x.SetReply(r)
		x.Authoritative = r.Authoritative

		return err
	}

	return (*w).WriteMsg(x)
}

func (h *DNSFSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()

	aType := r.Question[0].Qtype == dns.TypeA

	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true
	domain := msg.Question[0].Name

	if aType {
		if h.verbose {
			h.logger.Log("A; %v", domain)
		}

		if h.check(domain) {
			if h.verbose {
				h.logger.Log("(%v) sinking [%v from start]", domain, time.Since(start))
			}

			w.WriteMsg(&msg) // just sink right now
			return
		}
	} else if h.verbose {
		h.logger.Log("non A DNS request; forwarding")
	}

	go func() {
		if h.verbose {
			h.logger.Log("(%v) starting forwarding process [%v from start]", domain, time.Since(start))
		}

		for _, v := range h.forwards {
			if h.forward(&w, r, v) == nil {
				if h.verbose {
					h.logger.Log("(%v) passing on response from %v [%v from start]", domain, strings.Replace(v, ":", "#", 1), time.Since(start))
				}
				return
			}
		}

		h.ErrorChannel <- fmt.Errorf("errors upon all forwarding destinations for domain '%v'", domain)
	}()
}
