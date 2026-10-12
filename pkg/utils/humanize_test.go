package utils_test

import (
	"fmt"
	"testing"

	"github.com/nobbmaestro/lazyhis/pkg/utils"
)

func TestHumanizeDuration(t *testing.T) {
	cases := []struct {
		durationMs int64
		want       string
	}{
		{0, ""},
		{1, "  1ms"},
		{500, "500ms"},
		{999, "999ms"},
		{1000, "  1s"},
		{1500, "  1s"}, // integer division
		{59999, " 59s"},
		{60000, "  1m"},
		{3599999, " 59m"},
		{3600000, "  1h"},
		{86399999, " 23h"},
		{86400000, "  1d"},
		{172800000, "  2d"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%dms", tc.durationMs), func(t *testing.T) {
			got := utils.HumanizeDuration(tc.durationMs)

			if got != tc.want {
				t.Errorf(
					"HumanizeDuration(%d) = %q, want %q",
					tc.durationMs,
					got,
					tc.want,
				)
			}
		})
	}
}
