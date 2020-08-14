package models

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jcfox412/logarhythms/internal/utils"
)

const (
	headerPadding      = 2
	defaultTrackLength = 10 * time.Second
)

// Track is an object which can be played.
type Track struct {
	// Title of the track
	Title string
	// Amount of time the track should play
	Length time.Duration
	// Beats per minute (BPM) of the track
	BeatsPerMinute int
	// Number of beats per measure
	BeatsPerMeasure int
	// Number of times to devide each beat (e.g. 1 for quarter notes, 2 for eigth notes)
	DivisionsPerBeat int
	// Instruments used in the track.
	Instruments []*Instrument
	// Sequence of instruments to be played in the track.
	Patterns [][]*Instrument
}

// NewTrack creates a new track with calculated track pattern.
func NewTrack(title string, instruments []*Instrument, beatsPerMinute, beatsPerMeasure, divisionsPerBeat int) (*Track, error) {
	for _, instrument := range instruments {
		if err := instrument.validate(); err != nil {
			return nil, errors.Wrap(err, "error validating instruments")
		}
	}

	if err := validatePositiveInputs(beatsPerMinute, beatsPerMeasure, divisionsPerBeat); err != nil {
		return nil, errors.Wrap(err, "error validating integer inputs")
	}

	return &Track{
		Title:            title,
		Length:           defaultTrackLength,
		BeatsPerMinute:   beatsPerMinute,
		BeatsPerMeasure:  beatsPerMeasure,
		DivisionsPerBeat: divisionsPerBeat,
		Instruments:      instruments,
		Patterns:         makePattern(beatsPerMeasure*divisionsPerBeat, instruments),
	}, nil
}

// Play plays a track. This entails printing out the track's pattern as it is
// played, and playing the audio for the instruments of the track.
func (t *Track) Play() error {
	// delay allows for cleaner audio
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("Playing track at BPM: %v\n\n", t.BeatsPerMinute)

	header, headerWidth := t.printHeaders()
	fmt.Print(header)

	beatDuration, err := t.calculateBeatDuration()
	if err != nil {
		return errors.Wrap(err, "error calculating beat duration")
	}

	beatTicker := time.NewTicker(beatDuration)
	done := make(chan bool)

	go func() {
		beatDivisionCount := 0
		fmt.Print(utils.ClearLine(headerWidth))

		for {
			select {
			case <-done:
				return
			case <-beatTicker.C:
				if beatDivisionCount == t.BeatsPerMeasure*t.DivisionsPerBeat {
					beatDivisionCount = 0
					fmt.Print(utils.ClearLine(headerWidth))
				}

				fmt.Print(utils.CursorToNextColumn(len(t.Instruments) + 1))

				beatStr, err := t.triggerBeat(beatDivisionCount)
				if err != nil {
					// not great way to surface this error (since we're in a goroutine)
					fmt.Print(err.Error())
					return
				}

				fmt.Print(beatStr)

				beatDivisionCount++
			}
		}
	}()

	time.Sleep(t.Length)

	beatTicker.Stop()
	done <- true

	fmt.Println()

	return nil
}

func (t *Track) triggerBeat(beatDivisionCount int) (string, error) {
	beatStr := ""

	beatCount, err := utils.BeatCount(beatDivisionCount, t.DivisionsPerBeat)
	if err != nil {
		return "", errors.Wrap(err, "error determining beat count")
	}

	beatStr += beatCount
	beatStr += utils.CursorToNextRow()

	if beatDivisionCount >= len(t.Patterns) {
		return "", errors.New("beat counter higher than length of patterns - something went wrong")
	}

	instruments := t.Patterns[beatDivisionCount]

	for _, instrument := range instruments {
		if instrument != nil {
			instrument.Audio.Play()
			beatStr += fmt.Sprint("X|")

		} else {
			beatStr += fmt.Sprint("_|")
		}

		beatStr += utils.CursorToNextRow()
	}

	beatStr += utils.BeatTracker()

	return beatStr, nil
}

func (t *Track) calculateBeatDuration() (time.Duration, error) {
	if t.BeatsPerMinute <= 0 {
		return time.Duration(0), errors.New("beats per minute must be greater than 0")
	}

	if t.DivisionsPerBeat <= 0 {
		return time.Duration(0), errors.New("divisions per beat must be greater than 0")
	}

	beatDuration := 60.0 * 1000 / (t.BeatsPerMinute * t.DivisionsPerBeat)

	return time.Duration(beatDuration) * time.Millisecond, nil
}

func (t *Track) printHeaders() (string, int) {
	longestInstrument := 0
	header := ""

	for _, instrument := range t.Instruments {
		if len(instrument.Name) > longestInstrument {
			longestInstrument = len(instrument.Name)
		}
	}

	for _, instrument := range t.Instruments {
		header += fmt.Sprintf("%*s: |\n", longestInstrument, instrument.Name)
	}

	return header, longestInstrument + headerPadding
}

func validatePositiveInputs(beatsPerMinute, beatsPerMeasure, divisionsPerBeat int) error {
	if beatsPerMinute <= 0 {
		return errors.New("BeatsPerMinute must be greater than 0")
	}

	if beatsPerMeasure <= 0 {
		return errors.New("BeatsPerMeasure must be greater than 0")
	}

	if divisionsPerBeat <= 0 {
		return errors.New("DivisionsPerBeat must be greater than 0")
	}

	return nil
}

func makePattern(divisionsPerMeasure int, instruments []*Instrument) [][]*Instrument {
	pattern := make([][]*Instrument, divisionsPerMeasure)
	for i := 0; i < divisionsPerMeasure; i++ {
		pattern[i] = make([]*Instrument, len(instruments))
	}

	for i, instrument := range instruments {
		for _, beat := range instrument.Pattern {
			pattern[beat][i] = instrument
		}
	}

	return pattern
}
