package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToScaledVolume(t *testing.T) {
	type testCase struct {
		description    string
		input          float64
		expectedOutput float64
	}

	testCases := []testCase{
		{
			description:    "scales -5 to 0",
			input:          -5,
			expectedOutput: 0,
		},
		{
			description:    "scales 5 to 100",
			input:          5,
			expectedOutput: 100,
		},
		{
			description:    "scales 0 to 50",
			input:          0,
			expectedOutput: 50,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := toScaledVolume(testCase.input)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

func TestToUnscaledVolume(t *testing.T) {
	type testCase struct {
		description    string
		input          float64
		expectedOutput float64
	}

	testCases := []testCase{
		{
			description:    "scales 0 to -5",
			input:          0,
			expectedOutput: -5,
		},
		{
			description:    "scales 100 to 5",
			input:          100,
			expectedOutput: 5,
		},
		{
			description:    "scales 50 to 0",
			input:          50,
			expectedOutput: 0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := toUnscaledVolume(testCase.input)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}
