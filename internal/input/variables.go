package input

const (
	mainMenuOptions = "\n" +
		"1) Play Four on the floor (a pattern in 4/4 time)\n" +
		"2) Play Gravity by John Mayer (a pattern in 6/8 time)\n" +
		"3) Play Take Five by Dave Brubeck (a pattern in 5/4 time)\n" +
		"4) Experimental mode!\n" +
		"5) Exit\n"

	settingsMenuOptions = "\n" +
		"1) Beats per minute (BPM)\n" +
		"2) Instrument volume(s)\n" +
		"3) Track length\n" +
		"4) I'm done, play track!\n" +
		"5) Back to main menu\n"
)

var (
	inputTrackMap = map[string]string{
		"1": "assets/tracks/four_on_the_floor.json",
		"2": "assets/tracks/gravity.json",
		"3": "assets/tracks/take_five.json",
	}
)
