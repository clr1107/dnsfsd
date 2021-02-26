package rules

import (
	"regexp"
	"testing"
)

func TestWhitelist(t *testing.T) {
	rule := &containsRule{"google.com", false}
	whitelist := &containsRule{"456.google.com", true}

	set := &RuleSet{&map[IRule]struct{}{rule: {}, whitelist: {}}}

	results := [...]bool{
		set.Test("xxx.google.com"),
		set.Test("google.com"),
		set.Test("456.google.com"),
	}
	expected := [...]bool{true, true, false}

	if results != expected {
		t.Fatalf("incorrect results, received %v expected %v", results, expected)
	}
}

func TestRegexp(t *testing.T) {
	rule := &regexpRule{regexp.MustCompile(".*\\.google\\.com"), false}

	if !rule.Match("456.google.com") {
		t.Fatalf("regular expression rule did not match")
	}
}
