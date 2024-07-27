package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{Pending, "Pending"},
		{Scheduled, "Scheduled"},
		{Running, "Running"},
		{Completed, "Completed"},
		{Failed, "Failed"},
		{State(-1), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected State
		err      error
	}{
		{"Pending", Pending, nil},
		{"Scheduled", Scheduled, nil},
		{"Running", Running, nil},
		{"Completed", Completed, nil},
		{"Failed", Failed, nil},
		{"Unknown", State(-1), ErrStateUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			parsedState, err := Parse(tt.input)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, parsedState)
			}
		})
	}
}

func TestValidStateTransition(t *testing.T) {
	tests := []struct {
		src      State
		dst      State
		expected bool
	}{
		// Pending state transitions
		{src: Pending, dst: Scheduled, expected: true},
		{src: Pending, dst: Running, expected: false},
		{src: Pending, dst: Failed, expected: false},
		{src: Pending, dst: Completed, expected: false},

		// Scheduled state transitions
		{src: Scheduled, dst: Scheduled, expected: true},
		{src: Scheduled, dst: Running, expected: true},
		{src: Scheduled, dst: Failed, expected: true},
		{src: Scheduled, dst: Completed, expected: false},

		// Running state transitions
		{src: Running, dst: Running, expected: true},
		{src: Running, dst: Completed, expected: true},
		{src: Running, dst: Failed, expected: true},
		{src: Running, dst: Scheduled, expected: false},

		// Completed state transitions
		{src: Completed, dst: Pending, expected: false},
		{src: Completed, dst: Scheduled, expected: false},
		{src: Completed, dst: Running, expected: false},
		{src: Completed, dst: Completed, expected: false},
		{src: Completed, dst: Failed, expected: false},

		// Failed state transitions
		{src: Failed, dst: Pending, expected: false},
		{src: Failed, dst: Scheduled, expected: false},
		{src: Failed, dst: Running, expected: false},
		{src: Failed, dst: Completed, expected: false},
		{src: Failed, dst: Failed, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.src.String()+"_"+tt.dst.String(), func(t *testing.T) {
			result := ValidStateTransition(tt.src, tt.dst)
			assert.Equal(t, tt.expected, result)
		})
	}
}
