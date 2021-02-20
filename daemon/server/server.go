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
	return &handler{patterns, persistence.NewSimpleCache(-1)}
}

// true => sink; false => nothing found
func (h *handler) checkPatterns(domain string) bool {
	for _, pattern := range h.patterns {
		if pattern.MatchString(domain) {
			return true
		}
	}

	return false
}

// returns whether to sink or not based on cache and pattern matching
func (h *handler) check(domain string) bool {
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

func (h *handler) forward(w *dns.ResponseWriter, r *dns.Msg, msg *dns.Msg, dnsIP string) {
	// msg.Answer = append(msg.Answer, &dns.A{
	// 	Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
	// 	A:   net.ParseIP(result),
	// })
	(*w).WriteMsg(msg) // todo implement forwarding...
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	aType := r.Question[0].Qtype == dns.TypeA

	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	if aType {
		domain := msg.Question[0].Name

		if h.check(domain) {
			w.WriteMsg(&msg) // just sink right now
			return
		}
	}

	go h.forward(&w, r, &msg, "1.1.1.1")
}
