package agent

import (
	"github.com/stretchr/testify/assert"
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
			"app":     "Agent",
			"enabled": "true",
		}
		other := Selector{
			"app":                      "Agent",
			"enabled":                  "true",
			"some_other_key_we_ignore": "value",
		}
		assert.True(t, selector.Match(other))
	})

	t.Run("must_not_match_other_selector", func(t *testing.T) {
		selector := Selector{
			"app":     "Agent",
			"enabled": "false",
		}
		other := Selector{
			"app":                      "Agent",
			"enabled":                  "true",
			"some_other_key_we_ignore": "value",
		}
		assert.False(t, selector.Match(other))
	})

	t.Run("filterTargets_must_filter_out_non_matched_targets", func(t *testing.T) {
		targets := filterTargets(Targets{
			Target{
				Selector: Selector{
					"app":  "Agent",
					"name": "A",
				},
			},
			Target{
				Selector: Selector{
					"app":  "Agent",
					"name": "B",
				},
			},
			Target{
				Selector: Selector{
					"app":  "non-Agent",
					"name": "C",
				},
			},
			Target{
				Selector: Selector{
					"app":  "Agent",
					"name": "D",
				},
			},
		}, Selector{
			"app": "Agent",
		})
		assert.Len(t, targets, 3)
		assert.Equal(t, "A", targets[0].Selector["name"])
		assert.Equal(t, "B", targets[1].Selector["name"])
		assert.Equal(t, "D", targets[2].Selector["name"])
	})
}
