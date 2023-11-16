package templicons

import (
	"testing"
)

func TestIconURL(t *testing.T) {
	td := []struct{
		tn       string
		set      string
		name     string
		params   *Parameters
		expected string
	}{
		{
			set:      "mdi",
			name:     "home",
			params:   nil,
			expected: "/mdi/home.svg",
		},
		{
			set:      "mdi",
			name:     "home",
			params:   &Parameters{Color: "red"},
			expected: "/mdi/home.svg?color=red",
		},
	}

	for _, d := range td {
		t.Run(d.tn, func(t *testing.T) {
			got := iconPath(d.set, d.name, d.params)
			if got != d.expected {
				t.Errorf("expected %s, got %s", d.expected, got)
			}
		})
	}
}

