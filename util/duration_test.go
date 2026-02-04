package util

import (
	"strconv"
	"testing"

	"github.com/LeoRiether/gh-box/test/assert"
)

func TestParseSinceDays(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Duration
		wantErr error
	}{
		{name: "days suffix", input: "14d", want: 14 * Day},
		{name: "weeks suffix", input: "2w", want: 14 * Day},
		{name: "plain number", input: "7", want: 7 * Day},
		{name: "zero string", input: "0", want: 0},
		{name: "trimmed", input: " 3d ", want: 3 * Day},
		{name: "empty", input: "", want: 0},
		{name: "invalid suffix", input: "3x", wantErr: errInvalidUnit},
		{name: "negative", input: "-1", wantErr: errMustBePositive},
		{name: "nonnumeric", input: "abc", wantErr: &strconv.NumError{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if tt.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.ErrorAs(t, err, &tt.wantErr)
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
