package utils_test

import (
	"testing"

	"github.com/jcfox412/logarhythms/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestBeatCount(t *testing.T) {
	type input struct {
		beatCount        int
		divisionsPerBeat int
	}

	type testCase struct {
		description     string
		input           input
		expectedOutput  string
		expectedToError bool
	}

	testCases := []testCase{
		{
			description: "Succeeds on-beat",
			input: input{
				beatCount:        0,
				divisionsPerBeat: 2,
			},
			expectedOutput:  "1 ",
			expectedToError: false,
		},
		{
			description: "Succeeds off-beat",
			input: input{
				beatCount:        1,
				divisionsPerBeat: 2,
			},
			expectedOutput:  "  ",
			expectedToError: false,
		},
		{
			description: "Errors when divisionsPerBeat is 0",
			input: input{
				divisionsPerBeat: 0,
			},
			expectedOutput:  "",
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput, actualErr := utils.BeatCount(testCase.input.beatCount, testCase.input.divisionsPerBeat)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
			assert.Equal(t, testCase.expectedOutput, actualOutput)
		}

	}
}

func TestBold(t *testing.T) {
	type testCase struct {
		description    string
		input          string
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			input:          "boldme",
			expectedOutput: "\x1b[1mboldme\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := utils.Bold(testCase.input)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

func TestBeatTracker(t *testing.T) {
	type testCase struct {
		description    string
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			expectedOutput: "\x1b[3D   *",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := utils.BeatTracker()
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

func TestCursorToNextColumn(t *testing.T) {
	type testCase struct {
		description    string
		input          int
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			input:          4,
			expectedOutput: "\x1b[4A\x1b[1C",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := utils.CursorToNextColumn(testCase.input)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

func TestCursorToNextRow(t *testing.T) {
	type testCase struct {
		description    string
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			expectedOutput: "\x1b[1B\x1b[2D",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := utils.CursorToNextRow()
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

func TestClearLine(t *testing.T) {
	type testCase struct {
		description    string
		input          int
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			input:          10,
			expectedOutput: "\x1b[G\x1b[K\x1b[10C",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput := utils.ClearLine(testCase.input)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}
