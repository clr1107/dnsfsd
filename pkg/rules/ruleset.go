package rules

import (
	"regexp"
	"strings"
)

// RuleSet is a set of IRule implementations.
type RuleSet struct {
	rules map[IRule]struct{}
}

// Size returns the number of rules in this set.
func (s *RuleSet) Size() int {
	return len(s.rules)
}

// Test returns true if a given domain should be sinkholed. False indicates it
// should not. Whitelist rules are tested first as they always take precedence;
// any whitelist rule that matches will provide an immediate false indication.
// Blacklists are tested after. If there are no whitelist matches and no
// blacklist matches then no a false indication is given.
func (s *RuleSet) Test(domain string) bool {
	first := true
loop:
	for v := range s.rules {
		if (first == v.Whitelist()) && v.Match(domain) {
			return !first
		}
	}
	if first {
		first = false
		goto loop
	}

	return false
}

func ruleToString(prefix string, str string, whitelist bool) string {
	s := prefix + ";"

	if whitelist {
		s += "w"
	}

	return s + ";" + str
}

// IRule is an interface for a domain matching rule. Match returns true is the
// rule matches a given domain (case insensitive) Whitelist returns true if this
// is a whitelist rule; false if it is a blacklist rule String returns a string
// representation.
type IRule interface {
	Match(domain string) bool
	Whitelist() bool
	String() string
}

type regexpRule struct {
	expression *regexp.Regexp
	whitelist  bool
}

func (r regexpRule) Match(domain string) bool {
	return r.expression.MatchString(domain)
}

func (r regexpRule) Whitelist() bool {
	return r.whitelist
}

func (r regexpRule) String() string {
	return ruleToString(regexpRulePrefix, r.expression.String(), r.whitelist)
}

type containsRule struct {
	substring string
	whitelist bool
}

func (r containsRule) Match(domain string) bool {
	return strings.Contains(domain, r.substring)
}

func (r containsRule) Whitelist() bool {
	return r.whitelist
}

func (r containsRule) String() string {
	return ruleToString(containsRulePrefix, r.substring, r.whitelist)
}

type equalsRule struct {
	str       string
	whitelist bool
}

func (e equalsRule) Match(domain string) bool {
	return len(domain) == len(e.str) && domain == e.str
}

func (e equalsRule) Whitelist() bool {
	return e.whitelist
}

func (e equalsRule) String() string {
	return ruleToString(equalsRulePrefix, e.str, e.whitelist)
}
