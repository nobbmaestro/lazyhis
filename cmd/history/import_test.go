package history

import (
	"reflect"
	"testing"
)

func TestParseHistoryLine(t *testing.T) {
	cases := []struct {
		description string
		line        string
		want        []string
	}{
		{
			"empty line returns empty slice",
			"",
			[]string{},
		},
		{
			"plain single word",
			"ls",
			[]string{"ls"},
		},
		{
			"plain multi-word command",
			"git status",
			[]string{"git", "status"},
		},
		{
			"plain command with flags",
			"ls -la /tmp",
			[]string{"ls", "-la", "/tmp"},
		},
		{
			"zsh extended format",
			": 1234567890:0;git status",
			[]string{"git", "status"},
		},
		{
			"zsh extended format with duration",
			": 1234567890:42;docker run -it ubuntu",
			[]string{"docker", "run", "-it", "ubuntu"},
		},
		{
			"zsh extended format with empty command",
			": 1234567890:0;",
			[]string{},
		},
		{
			"zsh prefix without semicolon returns empty",
			": 1234567890:0",
			[]string{},
		},
		{
			"non-printable chars are stripped, words stay separate",
			"git\x01 status",
			[]string{"git", "status"},
		},
		{
			"only whitespace returns empty slice",
			"   ",
			[]string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			got := parseHistoryLine(tc.line)

			if len(tc.want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseHistoryLine(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}
