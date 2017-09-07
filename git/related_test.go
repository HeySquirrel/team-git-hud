package git

import (
	"testing"
)

func TestRelatedWorkItems(t *testing.T) {
	cases := []struct {
		Subject  string
		Expected []string
	}{
		{"F198234_team_coolness- scope coolness by user", []string{"F198234"}},
		{"S9028: Make something cool", []string{"S9028"}},
		{"DE9283: user is uncool", []string{"DE9283"}},
		{"S28973: Remove F2938_uncool_users", []string{"S28973", "F2938"}},
		{"No related work", []string{}},
	}

	for _, c := range cases {
		entry := new(LogEntry)
		entry.Subject = c.Subject

		entries := Logs{entry}
		actual := entries.relatedWorkItems()

		if len(actual) != len(c.Expected) {
			t.Fatalf("'%v' not equal '%v'", actual, c.Expected)
		}

		for i := range actual {
			if actual[i] != c.Expected[i] {
				t.Errorf("'%s' not equal '%s'", actual[i], c.Expected[i])
			}
		}
	}
}