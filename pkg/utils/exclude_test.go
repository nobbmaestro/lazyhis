package utils_test

import (
	"testing"

	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

func TestIsExcluded(t *testing.T) {
	cases := []struct {
		description string
		command     string
		exclusions  []string
		expected    bool
	}{
		{"exact match", "foo", []string{"^foo$"}, true},
		{"partial match should not be excluded", "foo bar", []string{"^foo$"}, false},
		{"prefix match should be excluded", "foo bar", []string{"^foo"}, true},
		{"contains match should be excluded", "foo bar", []string{"bar"}, true},
		{"regex match should be excluded", "foobar", []string{"foo.*"}, true},
		{"empty exclusions should not exclude", "foo", []string{}, false},
		{"empty command should not be excluded", "", []string{"^foo$"}, false},
		{"empty string exclusion should exclude empty command", "", []string{""}, true},
		{"multiple exclusions", "baz", []string{"foo", "bar", "baz"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			result := utils.MatchesExclusionPatterns(tc.command, tc.exclusions)
			if result != tc.expected {
				t.Errorf("Expected %q with exclusions %v to return %v, but got %v",
					tc.command, tc.exclusions, tc.expected, result)
			}
		})
	}
}
