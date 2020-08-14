package models_test

import (
	"testing"
	"time"

	"github.com/jcfox412/logarhythms/internal/audio"
	audiomocks "github.com/jcfox412/logarhythms/internal/audio/mocks"
	"github.com/jcfox412/logarhythms/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewTrack(t *testing.T) {
	type input struct {
		instruments      []*models.Instrument
		beatsPerMinute   int
		beatsPerMeasure  int
		divisionsPerBeat int
	}

	type testCase struct {
		description     string
		input           input
		expectedToError bool
	}

	testCases := []testCase{
		{
			description: "Successfully initializes Track object",
			input: input{
				instruments: []*models.Instrument{
					{Audio: &audio.BeepManager{}},
				},
				beatsPerMinute:   120,
				beatsPerMeasure:  4,
				divisionsPerBeat: 2,
			},
			expectedToError: false,
		},
		{
			description: "Fails with invalid instrument",
			input: input{
				instruments: []*models.Instrument{
					{},
				},
				beatsPerMinute:   120,
				beatsPerMeasure:  4,
				divisionsPerBeat: 2,
			},
			expectedToError: true,
		},
		{
			description: "Fails with invalid beatsPerMinute",
			input: input{
				instruments: []*models.Instrument{
					{Audio: &audio.BeepManager{}},
				},
				beatsPerMinute:   0,
				beatsPerMeasure:  4,
				divisionsPerBeat: 2,
			},
			expectedToError: true,
		},
		{
			description: "Fails with invalid beatsPerMeasure",
			input: input{
				instruments: []*models.Instrument{
					{Audio: &audio.BeepManager{}},
				},
				beatsPerMinute:   120,
				beatsPerMeasure:  0,
				divisionsPerBeat: 2,
			},
			expectedToError: true,
		},
		{
			description: "Fails with invalid divisionsPerBeat",
			input: input{
				instruments: []*models.Instrument{
					{Audio: &audio.BeepManager{}},
				},
				beatsPerMinute:   120,
				beatsPerMeasure:  4,
				divisionsPerBeat: 0,
			},
			expectedToError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		actualOutput, actualErr := models.NewTrack("testTitle", testCase.input.instruments, testCase.input.beatsPerMinute, testCase.input.beatsPerMeasure, testCase.input.divisionsPerBeat)
		if testCase.expectedToError {
			assert.Nil(t, actualOutput)
			assert.NotNil(t, actualErr)
		} else {
			assert.NotNil(t, actualOutput)
			assert.Nil(t, actualErr)
		}

	}
}

// TODO: this test does not cover the Play functionality super effectively. If
// I spent more time on this, I would change Play to make it more testable.
func TestPlay(t *testing.T) {
	type input struct {
		instruments      []*models.Instrument
		beatsPerMinute   int
		beatsPerMeasure  int
		divisionsPerBeat int
	}

	type testCase struct {
		description     string
		input           *models.Track
		setupMocks      func(*audiomocks.Manager)
		expectedToError bool
	}

	testInstrument1 := &models.Instrument{Name: "testInstrument1", Pattern: []int{0, 1, 2}}
	testInstrument2 := &models.Instrument{Name: "testInstrument2", Pattern: []int{1, 2}}

	testCases := []testCase{
		{
			description: "Successfully plays track",
			input: &models.Track{
				Length: 100 * time.Millisecond,
				Instruments: []*models.Instrument{
					testInstrument1,
					testInstrument2,
				},
				Patterns: [][]*models.Instrument{
					[]*models.Instrument{
						testInstrument1,
						testInstrument2,
					},
					[]*models.Instrument{},
				},
				BeatsPerMinute:   1000,
				BeatsPerMeasure:  2,
				DivisionsPerBeat: 1,
			},
			setupMocks: func(m *audiomocks.Manager) {
				m.On("Play").Return().Once()
			},
			expectedToError: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		track := testCase.input

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

		actualErr := track.Play()
		if testCase.expectedToError {
			assert.NotNil(t, actualErr)
		} else {
			assert.Nil(t, actualErr)
		}

		for _, m := range mockManagers {
			m.AssertExpectations(t)
		}
	}
}
