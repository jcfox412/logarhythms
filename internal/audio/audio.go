package audio

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
)

// Manager allows for mocking, which allows for easier testing of other packages.
type Manager interface {
	GetVolume() float64
	SetVolume(float64) (float64, error)
	Play()
}

// BeepManager manages audio state and functionality using the beep library.
type BeepManager struct {
	// Volume of the audio object. Valid between -5 and 5. See
	// https://godoc.org/github.com/faiface/beep/effects#Volume for full details
	// of how Volume works, but essentially input signal is multiplied by
	// math.Pow(2, Volume).
	volume float64
	// Buffer of audio data so file doesn't need to be opened every time it's played.
	buffer *beep.Buffer
}

var _ Manager = new(BeepManager)

// New creates a new audio Manager from the given wav filename. Returns an error
// if the file cannot be found or is not decodeable in the wav format.
func New(wavFilename string) (Manager, error) {
	buffer, err := setupSound(wavFilename)
	if err != nil {
		return nil, err
	}

	return &BeepManager{
		volume: 0,
		buffer: buffer,
	}, nil
}

// GetVolume fetches the Manager's scaled volume, from 0 to 100.
func (m *BeepManager) GetVolume() float64 {
	return toScaledVolume(m.volume)
}

// SetVolume sets the Manager's volume. Accepts values between 0 and 100.
func (m *BeepManager) SetVolume(volume float64) (float64, error) {
	if volume < 0 || volume > 100 {
		return volume, errors.New("volume must be between 0 and 100")
	}

	m.volume = toUnscaledVolume(volume)

	return toScaledVolume(m.volume), nil
}

// Play triggers audio to be played through the computer's default speaker.
// NOTE: this function is not tested, as adding testing would be more hassle
// than it's worth for a homework project.
func (m *BeepManager) Play() {
	speaker.Play(&effects.Volume{
		Streamer: m.buffer.Streamer(0, m.buffer.Len()),
		Base:     2,
		Volume:   m.volume,
		Silent:   false,
	})
}

func setupSound(wavFilename string) (*beep.Buffer, error) {
	f, err := os.Open(wavFilename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening sound file")
	}

	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding sound file")
	}

	// 40 found to sound best through experimentation
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/40))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	streamer.Close()

	return buffer, nil
}

func toScaledVolume(volume float64) float64 {
	return (volume + 5) * 10
}

func toUnscaledVolume(volume float64) float64 {
	return (volume / 10) - 5
}
