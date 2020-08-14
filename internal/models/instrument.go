package models

import (
	"github.com/pkg/errors"

	"github.com/jcfox412/logarhythms/internal/audio"
)

// Instrument stores audio for later playback.
type Instrument struct {
	// Name of the instrument, e.g. Snare
	Name string
	// File location of the instrument's audio sample (relative to root of project)
	// Filename string
	// Beat subdivisions where the instrument should be triggered
	Pattern []int
	// Manager for audio of instrument
	Audio audio.Manager
}

// NewInstrument builds an Instrument object with Audio support.
func NewInstrument(name, filename string, pattern []int) (*Instrument, error) {
	audioManager, err := audio.New(filename)
	if err != nil {
		return nil, err
	}

	return &Instrument{
		Name:    name,
		Pattern: pattern,
		Audio:   audioManager,
	}, nil
}

func (i *Instrument) validate() error {
	if i.Audio == nil {
		return errors.New("instrument audio manager must not be nil")
	}

	return nil
}
