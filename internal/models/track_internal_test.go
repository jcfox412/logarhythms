package models

import (
	"testing"
	"time"

	audiomocks "github.com/jcfox412/logarhythms/internal/audio/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPrintHeaders(t *testing.T) {
	type output struct {
		header      string
		headerWidth int
	}

	type testCase struct {
		description    string
		input          *Track
		expectedOutput output
	}

	testCases := []testCase{
		{
			description: "Handles no instruments",
			input:       &Track{},
			expectedOutput: output{
				header:      "",
				headerWidth: 2,
			},
		},
		{
			description: "Handles one instrument",
			input: &Track{
				Instruments: []*Instrument{
					{Name: "Snare"},
				},
			},
			expectedOutput: output{
				header:      "Snare: |\n",
				headerWidth: 7,
			},
		},
		{
			description: "Handles multiple instrument",
			input: &Track{
				Instruments: []*Instrument{
					{Name: "Kick"},
					{Name: "Snare"},
				},
			},
			expectedOutput: output{
				header:      " Kick: |\nSnare: |\n",
				headerWidth: 7,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualHeader, actualHeaderWidth := testCase.input.printHeaders()
		assert.Equal(t, testCase.expectedOutput.header, actualHeader)
		assert.Equal(t, testCase.expectedOutput.headerWidth, actualHeaderWidth)
	}
}

func TestCalculateBeatDuration(t *testing.T) {
	type testCase struct {
		description     string
		input           *Track
		expectedOutput  time.Duration
		expectedToError bool
	}

	testCases := []testCase{
		{
			description:     "Errors when beats per minute is 0",
			input:           &Track{},
			expectedOutput:  time.Duration(0),
			expectedToError: true,
		},
		{
			description:     "Errors when divisions per beat is 0",
			input:           &Track{BeatsPerMinute: 1},
			expectedOutput:  time.Duration(0),
			expectedToError: true,
		},
		{
			description:     "Succeeds",
			input:           &Track{BeatsPerMinute: 60, DivisionsPerBeat: 1},
			expectedOutput:  time.Duration(1000) * time.Millisecond,
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualDuration, actualErr := testCase.input.calculateBeatDuration()
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
			assert.Equal(t, testCase.expectedOutput, actualDuration)
		}
	}
}

func TestMakePattern(t *testing.T) {
	type input struct {
		divisionsPerMeasure int
		instruments         []*Instrument
	}

	type testCase struct {
		description    string
		input          input
		expectedOutput [][]*Instrument
	}

	testInstrument1 := &Instrument{Name: "testInstrument1", Pattern: []int{0, 1, 2}}
	testInstrument2 := &Instrument{Name: "testInstrument2", Pattern: []int{1, 2}}
	testPattern := [][]*Instrument{
		{testInstrument1, nil},
		{testInstrument1, testInstrument2},
		{testInstrument1, testInstrument2},
		{nil, nil},
	}

	testCases := []testCase{
		{
			description: "Succeeds with no instruments",
			input: input{
				divisionsPerMeasure: 0,
				instruments:         []*Instrument{},
			},
			expectedOutput: [][]*Instrument{},
		},
		{
			description: "Succeeds with instruments",
			input: input{
				divisionsPerMeasure: 4,
				instruments: []*Instrument{
					testInstrument1,
					testInstrument2,
				},
			},
			expectedOutput: testPattern,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualPattern := makePattern(testCase.input.divisionsPerMeasure, testCase.input.instruments)
		assert.Equal(t, len(testCase.expectedOutput), len(actualPattern))
		assert.Equal(t, testCase.expectedOutput, actualPattern)
	}
}

func TestTriggerBeat(t *testing.T) {
	type input struct {
		track     *Track
		beatCount int
	}

	type testCase struct {
		description     string
		input           input
		setupMocks      func(*audiomocks.Manager)
		expectedOutput  string
		expectedToError bool
	}

	testInstrument1 := &Instrument{Name: "testInstrument1"}
	testInstrument2 := &Instrument{Name: "testInstrument2"}

	testCases := []testCase{
		{
			description: "Errors with divisionsPerBeat is 0",
			input: input{
				track:     &Track{},
				beatCount: 0,
			},
			setupMocks:      func(m *audiomocks.Manager) {},
			expectedOutput:  "",
			expectedToError: true,
		},
		{
			description: "Succeeds with no instruments",
			input: input{
				track: &Track{
					Patterns: [][]*Instrument{
						{},
					},
					DivisionsPerBeat: 1,
				},
				beatCount: 0,
			},
			setupMocks:      func(m *audiomocks.Manager) {},
			expectedOutput:  "1 \x1b[1B\x1b[2D\x1b[3D   *",
			expectedToError: false,
		},
		{
			description: "Succeeds with one instrument playing and one not",
			input: input{
				track: &Track{
					Instruments: []*Instrument{
						testInstrument1,
						testInstrument2,
					},
					Patterns: [][]*Instrument{
						{
							testInstrument1,
							nil,
						},
					},
					DivisionsPerBeat: 1,
				},
				beatCount: 0,
			},
			setupMocks: func(m *audiomocks.Manager) {
				m.On("Play").Return().Once()
			},
			expectedOutput:  "1 \x1b[1B\x1b[2DX|\x1b[1B\x1b[2D_|\x1b[1B\x1b[2D\x1b[3D   *",
			expectedToError: false,
		},
		{
			description: "Succeeds with two instruments playing",
			input: input{
				track: &Track{
					Instruments: []*Instrument{
						testInstrument1,
						testInstrument2,
					},
					Patterns: [][]*Instrument{
						{
							testInstrument1,
							testInstrument2,
						},
					},
					DivisionsPerBeat: 1,
				},
				beatCount: 0,
			},
			setupMocks: func(m *audiomocks.Manager) {
				m.On("Play").Return().Once()
			},
			expectedOutput:  "1 \x1b[1B\x1b[2DX|\x1b[1B\x1b[2DX|\x1b[1B\x1b[2D\x1b[3D   *",
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		track := testCase.input.track

		mockManagers := make([]*audiomocks.Manager, 0, len(track.Instruments))

		for i, instrument := range track.Instruments {
			m := &audiomocks.Manager{}

			// only mock play for instruments whose pattern says they should play
			if len(track.Patterns) > 0 && track.Patterns[0][i] != nil {
				testCase.setupMocks(m)
			}

			instrument.Audio = m
			mockManagers = append(mockManagers, m)
		}

		actualBeatStr, actualErr := track.triggerBeat(testCase.input.beatCount)
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
			assert.Equal(t, testCase.expectedOutput, actualBeatStr)
		}

		for _, m := range mockManagers {
			m.AssertExpectations(t)
		}
	}
}
