package input_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	audiomocks "github.com/jcfox412/logarhythms/internal/audio/mocks"
	"github.com/jcfox412/logarhythms/internal/input"
	"github.com/jcfox412/logarhythms/internal/models"
	_ "github.com/jcfox412/logarhythms/testing"
)

// TODO: test more than just exiting
func TestPrintMainMenu(t *testing.T) {
	type testCase struct {
		description     string
		input           []string
		expectedToError bool
	}

	testCases := []testCase{
		{
			description:     "Succeeds in exiting",
			input:           []string{"5"},
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		for _, i := range testCase.input {
			stdin.Write([]byte(fmt.Sprintf("%s\n", i)))
		}

		userInput := input.UserInput{
			Reader: &stdin,
		}

		actualErr := userInput.PrintMainMenu()
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
		}
	}
}

// TODO: test non-error functionality
func TestSettingsMenu(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		expectedToError bool
	}

	initialBeatsPerMinute := 100

	testCases := []testCase{
		{
			description:     "Errors on non-integer input",
			input:           "help",
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(fmt.Sprintf("%s\n", testCase.input)))

		userInput := input.UserInput{
			Reader: &stdin,
		}

		track := &models.Track{BeatsPerMinute: initialBeatsPerMinute}

		actualErr := userInput.PrintSettingsMenu(track)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
			assert.Equal(t, initialBeatsPerMinute, track.BeatsPerMinute)
		} else {
			assert.Nil(t, actualErr)

			userInputInt, _ := strconv.Atoi(testCase.input)
			assert.Equal(t, userInputInt, track.BeatsPerMinute)
		}
	}
}

func TestBeatsPerMinuteMenu(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		expectedToError bool
	}

	initialBeatsPerMinute := 100

	testCases := []testCase{
		{
			description:     "Errors on user input out of range",
			input:           "1001",
			expectedToError: true,
		},
		{
			description:     "Errors on non-integer input",
			input:           "help",
			expectedToError: true,
		},
		{
			description:     "Handles good input",
			input:           "999",
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(fmt.Sprintf("%s\n", testCase.input)))

		userInput := input.UserInput{
			Reader: &stdin,
		}

		track := &models.Track{BeatsPerMinute: initialBeatsPerMinute}

		actualErr := userInput.BeatsPerMinuteMenu(track)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
			assert.Equal(t, initialBeatsPerMinute, track.BeatsPerMinute)
		} else {
			assert.Nil(t, actualErr)

			userInputInt, _ := strconv.Atoi(testCase.input)
			assert.Equal(t, userInputInt, track.BeatsPerMinute)
		}
	}
}

// TODO: test more than error path
func TestAllInstrumentsVolumeMenu(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		setupMocks      func(*audiomocks.Manager)
		expectedToError bool
	}

	instrumentVolume := 50.0

	testCases := []testCase{
		{
			description: "Errors on user input out of range",
			input:       "0",
			setupMocks: func(m *audiomocks.Manager) {
				m.On("GetVolume").Return(instrumentVolume).Once()
			},
			expectedToError: true,
		},
		{
			description: "Errors on non-integer input",
			input:       "help",
			setupMocks: func(m *audiomocks.Manager) {
				m.On("GetVolume").Return(instrumentVolume).Once()
			},
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(fmt.Sprintf("%s\n", testCase.input)))

		userInput := input.UserInput{
			Reader: &stdin,
		}

		m := &audiomocks.Manager{}
		if testCase.setupMocks != nil {
			testCase.setupMocks(m)
		}

		instrument := &models.Instrument{Name: "testInstrument", Audio: m}

		track := &models.Track{Instruments: []*models.Instrument{instrument}}

		actualErr := userInput.AllInstrumentsVolumeMenu(track)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
		}
	}
}

func TestInstrumentVolumeMenu(t *testing.T) {
	type testCase struct {
		description        string
		input              string
		setupMocks         func(*audiomocks.Manager)
		expectedToError    bool
		shouldIncludeAudio bool
	}

	initialInstrumentVolume := 50.0
	updatedInstrumentVolume := 100.0

	testCases := []testCase{
		{
			description:        "Errors on user input out of range",
			input:              "101",
			shouldIncludeAudio: true,
			setupMocks: func(m *audiomocks.Manager) {
				m.On("GetVolume").Return(initialInstrumentVolume).Twice()
			},
			expectedToError: true,
		},
		{
			description:        "Errors on non-integer input",
			input:              "help",
			shouldIncludeAudio: true,
			setupMocks: func(m *audiomocks.Manager) {
				m.On("GetVolume").Return(initialInstrumentVolume).Twice()
			},
			expectedToError: true,
		},
		{
			description:        "Errors on nil instrument Audio",
			input:              "",
			shouldIncludeAudio: false,
			setupMocks:         func(m *audiomocks.Manager) {},
			expectedToError:    true,
		},
		{
			description:        "Handles good input",
			input:              "100",
			shouldIncludeAudio: true,
			setupMocks: func(m *audiomocks.Manager) {
				m.On("GetVolume").Return(initialInstrumentVolume).Once()
				m.On("SetVolume", updatedInstrumentVolume).Return(updatedInstrumentVolume, nil).Once()
				m.On("GetVolume").Return(updatedInstrumentVolume).Once()
			},
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(fmt.Sprintf("%s\n", testCase.input)))

		userInput := input.UserInput{
			Reader: &stdin,
		}

		m := &audiomocks.Manager{}
		if testCase.setupMocks != nil {
			testCase.setupMocks(m)
		}

		instrument := &models.Instrument{Audio: m}

		if !testCase.shouldIncludeAudio {
			instrument.Audio = nil
		}

		actualErr := userInput.InstrumentVolumeMenu(instrument)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)

			if testCase.shouldIncludeAudio {
				assert.Equal(t, initialInstrumentVolume, instrument.Audio.GetVolume())
			}
		} else {
			assert.Nil(t, actualErr)

			userInputFloat, _ := strconv.ParseFloat(testCase.input, 64)
			assert.Equal(t, userInputFloat, instrument.Audio.GetVolume())
		}

		m.AssertExpectations(t)
	}
}

func TestTrackLengthMenu(t *testing.T) {
	type testCase struct {
		description     string
		input           string
		expectedToError bool
	}

	initialTrackLength := 10

	testCases := []testCase{
		{
			description:     "Errors on user input out of range",
			input:           "0",
			expectedToError: true,
		},
		{
			description:     "Errors on non-integer input",
			input:           "help",
			expectedToError: true,
		},
		{
			description:     "Handles good input",
			input:           "1",
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		var stdin bytes.Buffer
		stdin.Write([]byte(fmt.Sprintf("%s\n", testCase.input)))

		userInput := input.UserInput{
			Reader: &stdin,
		}

		track := &models.Track{Length: time.Duration(initialTrackLength) * time.Second}

		actualErr := userInput.TrackLengthMenu(track)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
			assert.Equal(t, initialTrackLength, int(track.Length.Seconds()))
		} else {
			assert.Nil(t, actualErr)

			userInputInt, _ := strconv.Atoi(testCase.input)
			assert.Equal(t, userInputInt, int(track.Length.Seconds()))
		}
	}
}
