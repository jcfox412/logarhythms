package audio_test

import (
	"testing"

	"github.com/jcfox412/logarhythms/internal/audio"
	"github.com/stretchr/testify/assert"
)

func TestSetGetVolume(t *testing.T) {
	type testCase struct {
		description     string
		input           float64
		expectedOutput  float64
		expectedToError bool
	}

	testCases := []testCase{
		{
			description:     "Sets to minimum volume",
			input:           0,
			expectedOutput:  0,
			expectedToError: false,
		},
		{
			description:     "Sets to maximum volume",
			input:           100,
			expectedOutput:  100,
			expectedToError: false,
		},
		{
			description:     "Errors with volume too low",
			input:           -1,
			expectedToError: true,
		},
		{
			description:     "Errors with volume too low",
			input:           101,
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		a := audio.BeepManager{}

		setOutput, setErr := a.SetVolume(testCase.input)
		if testCase.expectedToError {
			assert.NotNil(t, setErr)

			getOutput := a.GetVolume()
			assert.Equal(t, 50.0, getOutput)
		} else {
			assert.Nil(t, setErr)
			assert.Equal(t, testCase.expectedOutput, setOutput)

			getOutput := a.GetVolume()
			assert.Equal(t, testCase.expectedOutput, getOutput)
		}

	}
}

func TestNew(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		expectedToError bool
	}

	testCases := []testCase{
		{
			description:     "Successfully initializes Audio object",
			input:           "testfiles/valid.wav",
			expectedToError: false,
		},
		{
			description:     "Fails with nonexistant audio file",
			input:           "testfiles/nonexistant.wav",
			expectedToError: true,
		},
		{
			description:     "Fails with bad audio file",
			input:           "testfiles/invalid.wav",
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput, actualErr := audio.New(testCase.input)
		if testCase.expectedToError {
			assert.Nil(t, actualOutput)
			assert.NotNil(t, actualErr)
		} else {
			assert.NotNil(t, actualOutput)
			assert.Nil(t, actualErr)
		}

	}
}
