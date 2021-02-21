package rules

import (
	"regexp"
	"strings"
)

type RuleSet struct {
	rules *[]IRule
}

func (s *RuleSet) Size() int {
	return len(*s.rules)
}

func (s *RuleSet) Test(domain string) bool {
	for _, v := range *s.rules {
		if v.Whitelist() {
			if v.Match(domain) {
				return false
			}
		}
	}

	for _, v := range *s.rules {
		if !v.Whitelist() {
			if v.Match(domain) {
				return true
			}
		}
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
	return r.expression.MatchString(strings.ToLower(domain))
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
	return strings.Contains(strings.ToLower(domain), r.substring)
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
	return (strings.ToLower(domain) == e.str)
}

func (e equalsRule) Whitelist() bool {
	return e.whitelist
}

func (e equalsRule) String() string {
	return ruleToString(equalsRulePrefix, e.str, e.whitelist)
}
