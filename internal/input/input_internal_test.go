package input

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/jcfox412/logarhythms/internal/models"
	_ "github.com/jcfox412/logarhythms/testing"
)

func TestPrepareTrack(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		expectedOutput  *models.Track
		expectedToError bool
	}

	testCases := []testCase{
		{
			description: "Successfully creates valid track",
			input:       "internal/input/testfiles/valid_track.json",
			expectedOutput: &models.Track{
				Instruments: []*models.Instrument{
					{},
				},
				Patterns:         make([][]*models.Instrument, 8),
				Title:            "Valid Track",
				BeatsPerMeasure:  4,
				DivisionsPerBeat: 2,
				BeatsPerMinute:   120,
			},
			expectedToError: false,
		},
		{
			description:     "Errors on nonexistant track file",
			input:           "internal/input/testfiles/nonexistant.json",
			expectedOutput:  &models.Track{},
			expectedToError: true,
		},
		{
			description:     "Errors on invalid track file",
			input:           "internal/input/testfiles/invalid_track.json",
			expectedOutput:  &models.Track{},
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput, actualErr := prepareTrack(testCase.input)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
			compareTracks(t, testCase.expectedOutput, actualOutput)
		}
	}
}

func TestRetry(t *testing.T) {
	type retryInput struct {
		attempts             int
		expectedCallAttempts int
		f                    func(interface{}) error
	}

	type testCase struct {
		description     string
		input           retryInput
		expectedToError bool
	}

	type callCounter struct {
		count int
	}

	testCases := []testCase{
		{
			description: "Successfully executes function",
			input: retryInput{
				attempts:             1,
				expectedCallAttempts: 1,
				f: func(i interface{}) error {
					callCount := i.(*callCounter)
					callCount.count++
					return nil
				},
			},
			expectedToError: false,
		},
		{
			description: "Succeeds within retry attempts",
			input: retryInput{
				attempts:             3,
				expectedCallAttempts: 2,
				f: func(i interface{}) error {
					callCount := i.(*callCounter)
					callCount.count++

					if callCount.count == 2 {
						return nil
					}

					return errors.New("halp")
				},
			},
			expectedToError: false,
		},
		{
			description: "Errors after appropriate number of attempts",
			input: retryInput{
				attempts:             3,
				expectedCallAttempts: 3,
				f: func(i interface{}) error {
					callCount := i.(*callCounter)
					callCount.count++
					return errors.New("halp")
				},
			},
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		callCount := &callCounter{count: 0}

		actualErr := retry(testCase.input.attempts, callCount, testCase.input.f)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
		}

		assert.Equal(t, testCase.input.expectedCallAttempts, callCount.count)
	}
}

func TestGetUserInput(t *testing.T) {
	type testCase struct {
		description    string
		input          string
		expectedOutput string
	}

	testCases := []testCase{
		{
			description:    "Succeeds",
			input:          "abc123\n",
			expectedOutput: "abc123",
		},
		{
			description:    "Succeeds with empty input",
			input:          "",
			expectedOutput: "",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(testCase.input))

		actualOutput := getUserInput(&stdin)
		assert.Equal(t, testCase.expectedOutput, actualOutput)
	}
}

// Helper function for comparing tracks without actually caring about audio
func compareTracks(t *testing.T, expectedTrack *models.Track, actualTrack *models.Track) {
	assert.Equal(t, expectedTrack.Title, actualTrack.Title)
	assert.Equal(t, expectedTrack.BeatsPerMeasure, actualTrack.BeatsPerMeasure)
	assert.Equal(t, expectedTrack.DivisionsPerBeat, actualTrack.DivisionsPerBeat)
	assert.Equal(t, expectedTrack.BeatsPerMinute, actualTrack.BeatsPerMinute)
	assert.Equal(t, len(expectedTrack.Instruments), len(actualTrack.Instruments))
	assert.Equal(t, len(expectedTrack.Patterns), len(actualTrack.Patterns))
}
