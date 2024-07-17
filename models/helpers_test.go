package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertToFixedTZ(t *testing.T) {
	testCases := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "UTC to PST",
			input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2019, 12, 31, 16, 0, 0, 0, time.FixedZone("PST", -8*60*60)),
		},
		{
			name:     "PST to PST",
			input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*60*60)),
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*60*60)),
		},
		{
			name:     "PDT to PST",
			input:    time.Date(2020, 1, 1, 1, 0, 0, 0, time.FixedZone("PDT", -7*60*60)),
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*60*60)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ConvertToFixedTZ(tc.input))
		})
	}
}
