package rules

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

func RuleFromString(text string) (IRule, error) {
	text = strings.TrimRight(text, "\n")
	text = strings.Trim(text, " ")

	if len(text) > 0 && text[0] == '#' {
		return nil, nil
	}

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
		return nil, fmt.Errorf("could not parse rule '%v' as it is in an invalid format", text)
	}

	switch split[0] {
	case regexpRulePrefix:
		pattern, err := regexp.Compile(ruleText)

		if err != nil {
			return nil, fmt.Errorf("could not parse rule as regular expression (opcode `r`) '%v'", text)
		}

		return regexpRule{pattern, whitelist}, nil
	case containsRulePrefix:
		return containsRule{ruleText, whitelist}, nil
	case equalsRulePrefix:
		return equalsRule{ruleText, whitelist}, nil
	default:
		return nil, fmt.Errorf("could not parse rule '%v' as opcode `%v` is unknown", text, split[0])
	}
}

// RuleFile is a representation of a file containing rules.
type RuleFile struct {
	Path   string   // Path to the file
	Loaded bool     // Whether the file has been loadaed yet
	Rules  *[]IRule // A pointer to a slice of rules that have been loaded
}

// Load loads a RuleFile and returns any errors.
func (p *RuleFile) Load() error {
	f, err := os.Open(p.Path)

	if err != nil {
		return fmt.Errorf("could not open rule file '%v'", p.Path)
	}

	scanner := bufio.NewScanner(f)
	rules := make([]IRule, 0)

	for scanner.Scan() {
		text := scanner.Text()
		rule, err := RuleFromString(text)

		if err != nil {
			return fmt.Errorf("%v: rule file %v", err, p.Path)
		}

		if rule != nil {
			rules = append(rules, rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	p.Rules = &rules
	p.Loaded = true

	return nil
}

// AllRulesFiles returns a pointer to a slice of RuleFiles inside a given
// directory, and any errors encountered whislst reading the directory.
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

// LoadAllRuleFiles returns a pointer to a slice of *loaded* RuleFiles in a
// given directory and any errors encountered whilst reading and loading the
// files.
func LoadAllRuleFiles(path string) (*[]RuleFile, error) {
	files, err := AllRulesFiles(path)

	if err != nil {
		return nil, err
	}

	successes := make([]RuleFile, 0, len(*files))

	for _, v := range *files {
		if err := v.Load(); err == nil {
			successes = append(successes, v)
		} else {
			return nil, err
		}
	}

	return &successes, nil
}

// CollectAllRules creates a RuleSet from a pointer to a slice of RuleFiles. All
// RuleFiles must already be loaded otherwise they will be skipped.
func CollectAllRules(files *[]RuleFile) *RuleSet {
	l := make(map[IRule]struct{})

	for _, v := range *files {
		if v.Loaded {
			for _, rule := range *v.Rules {
				l[rule] = struct{}{}
			}
		}
	}

	return &RuleSet{&l}
}

// DownloadRuleFile downloads over http from a given URL to /etc/dnsfsd/rules
// and a given file name. It returns the number of rules in the file and any
// errors encountered.
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
