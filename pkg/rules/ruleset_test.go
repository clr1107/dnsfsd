package rules

import (
	"regexp"
	"testing"
)

func TestWhitelist(t *testing.T) {
	rule := &containsRule{"google.com", false}
	whitelist := &containsRule{"456.google.com", true}

	set := &RuleSet{map[IRule]struct{}{rule: {}, whitelist: {}}}

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

func TestContains(t *testing.T) {
	rule := &containsRule{"substring", false}

	if !rule.Match("x.y.substring.google.com") {
		t.Fatalf("contains rule did not match")
	}
}

func TestEquals(t *testing.T) {
	rule := &equalsRule{"google.com", false}

	if !rule.Match("google.com") {
		t.Fatalf("equals rule did not match")
	}
}

func TestRuleFromString(t *testing.T) {
	var err error
	var one IRule
	var two IRule
	var three IRule
	errors := [3]error{}

	one, err = RuleFromString("r;w;regexhere")
	errors[0] = err

	two, err = RuleFromString("e;;stringhere")
	errors[1] = err

	three, err = RuleFromString("c;;substringhere")
	errors[2] = err

	for _, v := range errors {
		if v != nil {
			t.Fatalf("error ocurred: %v", v)
		}
	}

	if _, ok := one.(regexpRule); !ok {
		t.Fatalf("regex rule was not parsed as such")
	}

	if _, ok := two.(equalsRule); !ok {
		t.Fatalf("equals rule was not parsed as such")
	}

	if _, ok := three.(containsRule); !ok {
		t.Fatalf("contains rule was not parsed as such")
	}

	if _, err = RuleFromString("thisisinvalid"); err == nil {
		t.Fatal("no error for invalid rule")
	}
}
