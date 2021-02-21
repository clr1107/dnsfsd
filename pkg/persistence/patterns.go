package persistence

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	regexpRulePrefix   string = "r"
	containsRulePrefix string = "c"
	equalsRulePrefix   string = "e"
	whitelistChar      rune   = 'w'
)

type IRule interface {
	Match(domain string) bool
	String() string
}

func ruleToString(prefix string, str string, whitelist bool) string {
	s := prefix + ";"

	if whitelist {
		s += "w"
	}

	return s + ";" + str
}

type regexpRule struct {
	expression *regexp.Regexp
	whitelist  bool
}

func (r regexpRule) Match(domain string) bool {
	return r.expression.MatchString(strings.ToLower(domain)) != r.whitelist
}

func (r regexpRule) String() string {
	return ruleToString(regexpRulePrefix, r.expression.String(), r.whitelist)
}

type containsRule struct {
	substring string
	whitelist bool
}

func (r containsRule) Match(domain string) bool {
	return strings.Contains(strings.ToLower(domain), r.substring) != r.whitelist
}

func (r containsRule) String() string {
	return ruleToString(containsRulePrefix, r.substring, r.whitelist)
}

type equalsRule struct {
	str       string
	whitelist bool
}

func (e equalsRule) Match(domain string) bool {
	return (strings.ToLower(domain) == e.str) != e.whitelist
}

func (e equalsRule) String() string {
	return ruleToString(equalsRulePrefix, e.str, e.whitelist)
}

type RuleFile struct {
	Path   string
	Loaded bool
	Rules  *[]IRule
}

func (p *RuleFile) Load() error {
	f, err := os.Open(p.Path)

	if err != nil {
		return fmt.Errorf("could not open rule file '%v'", p.Path)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	rules := make([]IRule, 0)

	for scanner.Scan() {
		text := scanner.Text()
		split := strings.SplitN(text, ";", 3)

		var ruleText string
		var whitelist bool = false

		if len(split) == 3 {
			if len(split[1]) > 0 {
				whitelist = rune(split[1][0]) == whitelistChar
			}

			ruleText = split[2]
		} else if len(split) == 2 {
			ruleText = split[1]
		} else {
			return fmt.Errorf("could not parse rule '%v' as it is in an invalid format", text)
		}

		switch split[0] {
		case regexpRulePrefix:
			pattern, err := regexp.Compile(ruleText)

			if err != nil {
				return fmt.Errorf("could not parse rule as regular expression (opcode `r`) '%v'", text)
			}

			rules = append(rules, regexpRule{pattern, whitelist})
			break
		case containsRulePrefix:
			rules = append(rules, containsRule{ruleText, whitelist})
			break
		case equalsRulePrefix:
			rules = append(rules, equalsRule{ruleText, whitelist})
			break
		default:
			return fmt.Errorf("could not parse rule '%v' as opcode `%v` is unknown", text, split[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	p.Rules = &rules
	p.Loaded = true

	return nil
}

func AllRulesFiles(directory string) (*[]RuleFile, error) {
	files, err := ioutil.ReadDir(directory)
	paths := make([]RuleFile, 0)

	if err != nil {
		return nil, fmt.Errorf("could not read all pattern files in directory '%v'", directory)
	}

	for _, v := range files {
		if !v.IsDir() {
			paths = append(paths, RuleFile{path.Join(directory, v.Name()), false, nil})
		}
	}

	return &paths, nil
}

func LoadAllRuleFiles(path string) (*[]RuleFile, error) {
	files, err := AllRulesFiles(path)
	successes := make([]RuleFile, 0, len(*files))

	if err != nil {
		return nil, err
	}

	for _, v := range *files {
		if err := v.Load(); err == nil {
			successes = append(successes, v)
		} else {
			return nil, err
		}
	}

	return &successes, nil
}

func CollectAllRules(files *[]RuleFile) *[]IRule {
	l := make([]IRule, 0, len(*files))

	for _, v := range *files {
		if v.Loaded {
			l = append(l, *v.Rules...)
		}
	}

	return &l
}

func DownloadRuleFile(url string, filename string) (int, error) {
	resp, err := http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("could not download from url: '%v'", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP request to '%v' gave a %v status code", url, resp.StatusCode)
	}

	const directory string = "/etc/dnsfsd/rules"
	if err := os.MkdirAll(directory, os.FileMode(0755)); err != nil {
		return 0, err
	}

	filepath := path.Join(directory, filename)
	out, err := os.Create(filepath)

	if err != nil {
		return 0, fmt.Errorf("could not create file '%v'", filepath)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}

	ruleFile := RuleFile{filepath, false, nil}

	if err := ruleFile.Load(); err != nil {
		return 0, err
	}

	return len(*ruleFile.Rules), nil
}
