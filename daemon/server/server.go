package server

import (
	"fmt"
	"github.com/clr1107/dnsfsd/pkg/data/cache"
	"strconv"
	"strings"

	"github.com/clr1107/dnsfsd/daemon/logger"
	"github.com/clr1107/dnsfsd/pkg/rules"
	"github.com/miekg/dns"
)

func newMsgReply(m *dns.Msg, ans []dns.RR) *dns.Msg {
	r := new(dns.Msg)
	r.SetReply(m)

	r.Authoritative = m.Authoritative
	r.Answer = ans

	return r
}

func formatDomain(domain string) string {
	domain = strings.ToLower(domain)
	l := len(domain)

	if l > 1 && domain[l - 1] == '.' {
		return domain[:len(domain) - 1]
	}

	return domain
}

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
	s.Handler.sinkCache.Clear()
	s.Handler.dnsCache.Clean()

	if err := s.Handler.dnsCache.SerialiseToFile("/etc/dnsfsd/dns.cache"); err != nil {
		return err
	}

	close(s.Handler.ErrorChannel)
	return s.Server.Shutdown()
}

type DNSFSHandler struct {
	rules        *rules.RuleSet
	sinkCache    *cache.SimpleCache
	dnsCache     *cache.DNSCache
	forwards     []string
	ErrorChannel chan error
	verbose      bool
	logger       *logger.Logger
}

func NewHandler(rules *rules.RuleSet, dnsCache *cache.DNSCache, forwards []string, verbose bool, logger *logger.Logger) *DNSFSHandler {
	return &DNSFSHandler{
		rules,
		cache.NewSimpleCache(-1),
		dnsCache,
		forwards,
		make(chan error),
		verbose,
		logger,
	}
}

// returns whether to sink or not based on cache and rule matching
func (h *DNSFSHandler) check(domain string) bool {
	if h.sinkCache.Contains(domain) {
		if val, ok := h.sinkCache.Get(domain).(bool); ok {
			return val
		}

		h.sinkCache.Remove(domain) // for some reason not a bool?
	}

	if h.rules.Test(domain) {
		h.sinkCache.PutDefault(domain, true)
		return true
	}

	h.sinkCache.PutDefault(domain, false)
	return false
}

func (h *DNSFSHandler) resolve(r *dns.Msg) (*dns.Msg, error) {
	question := r.Question[0]

	if h.dnsCache.Contains(question.String()) { // todo -- cache non-string keys
		rr, ok := h.dnsCache.Get(question.String()).([]dns.RR)

		if ok {
			return newMsgReply(r, rr), nil
		}

		h.dnsCache.Remove(question.String())
	}

	for k, v := range h.forwards {
		if h.verbose {
			h.logger.Log("[forwarding-%v] %v -> %v", k+1, question.String(), v)
		}

		msg, err := h.forward(r, v)

		if err != nil {
			h.ErrorChannel <- err
		} else {
			h.dnsCache.PutDefault(question.String(), msg.Answer)
			return msg, nil
		}
	}

	return nil, fmt.Errorf("no given DNS servers returned a result for this query: `%v`", question.String())
}

func (h *DNSFSHandler) forward(r *dns.Msg, dnsAddress string) (*dns.Msg, error) {
	question := r.Question[0]
	c := new(dns.Client)
	x, _, err := c.Exchange(r, dnsAddress)

	if err != nil || x == nil {
		if err == nil {
			err = fmt.Errorf("after forwarding query `%v` to '%v' the message response was nil", question.String(), dnsAddress)
		}

		return nil, err
	}

	return x, nil
}

func (h *DNSFSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	question := r.Question[0]
	domain := formatDomain(question.Name)

	if h.check(domain) {
		if err := w.WriteMsg(newMsgReply(r, nil)); err != nil {
			h.ErrorChannel <- err
			return
		}

		if h.verbose {
			h.logger.Log("[sink] %v", question.String())
		}

		return
	}

	go func() {
		msg, err := h.resolve(r)

		if err == nil {
			err = w.WriteMsg(msg)
		}

		if err != nil {
			h.ErrorChannel <- err

			if h.verbose {
				h.logger.LogErr("no response sent to question (err) %v", question.String())
			}
		}
	}()
}
