package agent

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

// Targets represent Agent per-target settings
type Targets []Target

// Target represents target
type Target struct {
	Name     string        `yaml:"name"`
	URI      string        `yaml:"uri"`
	Interval time.Duration `yaml:"interval"`
	RegExp   string        `yaml:"regexp"`
	Selector Selector      `yaml:"selector"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Selector represents selector description
type Selector map[string]string

// Match check if a given selector matches
func (s Selector) Match(other Selector) bool {
	for name, value := range s {
		otherValue, otherFound := other[name]
		if !otherFound {
			return false
		}
		if otherValue != value {
			return false
		}
	}
	return true
}

// MustParseSelector helper to convert selector passed as a string to the actual data structure
func MustParseSelector(selector []string) Selector {
	parsed, err := parseSelector(selector)
	mustNotErr(err)
	return parsed
}

func parseSelector(selector []string) (Selector, error) {
	parsed := Selector{}
	for _, s := range selector {
		parts := strings.Split(s, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("unparsable selector `%s`", s)
		}
		parsed[parts[0]] = parts[1]
	}
	return parsed, nil
}

// MustReadTargets reads target file, panics on failure
func MustReadTargets(p string, selector Selector) Targets {
	var targets Targets
	data, err := ioutil.ReadFile(p)
	mustNotErr(err)
	mustNotErr(yaml.Unmarshal(data, &targets))
	return filterTargets(targets, selector)
}

func filterTargets(in Targets, selector Selector) Targets {
	filtered := make(Targets, 0, len(in))
	for _, t := range in {
		if !selector.Match(t.Selector) {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}
