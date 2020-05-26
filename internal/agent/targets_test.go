package agent

import (
	"github.com/likexian/gokit/assert"
	"testing"
)

func TestSelectors(t *testing.T) {
	t.Run("must_parse_selector_string", func(t *testing.T) {
		rawSelector := []string{"app=value", "enabled=false"}
		expected := Selector{
			"app":     "value",
			"enabled": "false",
		}
		selector, err := parseSelector(rawSelector)
		assert.Nil(t, err)
		assert.Equal(t, expected, selector)
	})

	t.Run("must_match_other_selector", func(t *testing.T) {
		selector := Selector{
			"app":     "agent",
			"enabled": "true",
		}
		other := Selector{
			"app":                      "agent",
			"enabled":                  "true",
			"some_other_key_we_ignore": "value",
		}
		assert.True(t, selector.Match(other))
	})

	t.Run("must_not_match_other_selector", func(t *testing.T) {
		selector := Selector{
			"app":     "agent",
			"enabled": "false",
		}
		other := Selector{
			"app":                      "agent",
			"enabled":                  "true",
			"some_other_key_we_ignore": "value",
		}
		assert.False(t, selector.Match(other))
	})

	t.Run("filterTargets_must_filter_out_non_matched_targets", func(t *testing.T) {
		targets := filterTargets(Targets{
			Target{
				Selector: Selector{
					"app":  "agent",
					"name": "A",
				},
			},
			Target{
				Selector: Selector{
					"app":  "agent",
					"name": "B",
				},
			},
			Target{
				Selector: Selector{
					"app":  "non-agent",
					"name": "C",
				},
			},
			Target{
				Selector: Selector{
					"app":  "agent",
					"name": "D",
				},
			},
		}, Selector{
			"app": "agent",
		})
		assert.Len(t, targets, 3)
		assert.Equal(t, "A", targets[0].Selector["name"])
		assert.Equal(t, "B", targets[1].Selector["name"])
		assert.Equal(t, "D", targets[2].Selector["name"])
	})
}
