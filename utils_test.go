package templicons

import (
	"testing"
)

var api = "https://api.iconify.design"

func TestIconURL(t *testing.T) {
	td := []struct{
		tn       string
		api      string
		set      string
		name     string
		params   *Parameters
		expected string
	}{
		{
			api:      api,
			set:      "mdi",
			name:     "home",
			params:   nil,
			expected: "https://api.iconify.design/mdi/home",
		},
		{
			api:      api,
			set:      "mdi",
			name:     "home",
			params:   &Parameters{Color: "red"},
			expected: "https://api.iconify.design/mdi/home?color=red",
		},
	}

	for _, d := range td {
		t.Run(d.tn, func(t *testing.T) {
			got := iconURL(d.api, d.set, d.name, d.params)
			if got != d.expected {
				t.Errorf("expected %s, got %s", d.expected, got)
			}
		})
	}
}

