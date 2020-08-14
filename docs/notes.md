## Features

### Clean Playback User Interface

I'm satisfied with the way my playback UI turned out. It's achieved by using [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code), which I had no idea existed before this project.

### Audio Support

I initially played around with [portaudio](https://github.com/gordonklaus/portaudio) for audio support in this project, however I ended up using the [Beep](https://github.com/faiface/beep) library. I chose this library because it had slightly more usage in Github, was being kept up-to-date (last commit about a month ago), and had great documentation.

### Configurability

If you spend a little time in my test-based User Interface, you'll see that multiple aspects of a track can be changed. In addition to different BPMs, the project supports different lengths of tracks, as well as volume control for each individual instrument.

## If I Had More Time

### Testing

I'm generally happy with how I tested the package, however I feel that tests around the `input` package could have been more thorough. In general, it seems difficult to test user input in go, however I still feel that I could have abstracted my menu functionality out a bit better.

### User Interface

Speaking of the user interface, I decided early on that I would not be connecting this project to a UI. This was a conscious decision - I'm not a frontend engineer, so I figured my time would be better-spent implementing more backend features (e.g. volume control) rather than building out a frontend.

### Multi-Platform Support

I've tested this project fairly thoroughly on my own machine, and it works pretty much the way I want it to. That being said, I was able to test in a limited capacity on other platforms (why does my housemate own a Windows laptop) and the results were ... less than stellar.

On Windows, the audio was much choppier than on Mac, and it turns out the generic Windows terminal does NOT support ANSI escape codes, so the text printout was absolute garbage.

I don't have access to a Linux machine that has speakers, so I tested the project in a docker container, where the text printout worked perfectly (minus the audio libraries spewing errors when trying to initialize).

### Dynamic Controls

While my project does allow the user to control multiple aspects of playback before a track begins playing, it does not allow for control during playback. I toyed around with this concept, however ran into a wall when trying to figure out how to add keybindings during playback. I'm sure this is possible, probably even without adding a larger dependency like [tcell](https://github.com/gdamore/tcell).

## Potential Extensions

### User interface

I feel like the first logical extension for this project would be building out a web-based user interface. The terminal, while flexible in its own right, really does pale in comparison with what modern javascript frameworks can accomplish.

Building out this UI would be paired with server-ifying the go project.

### Easier addition of new tracks

Right now, this project features 3 tracks - Four on the Floor, Gravity, and Take Five. I've already done the work the support adding new tracks via new JSON files (also helpful for an API), however the user menus are still hard-coded to display only those three tracks. A simple extension of this project would be to have been to abstract the menu away from these tracks, so that any JSON could just be dropped in an would be able to be played.

### Datastore

Right now tracks live only in memory. This is great for simplicity, but fails when the system needs to handle hundreds or thousands of tracks or instruments. Adding a datastore of any sort could help the system scale with increasing object numbers. 


## General Thoughts

### Previous Art

When I initially started thinking about this exercise, I did a quick git search to see if any previous applicants had uploaded their code. As it turns out, I found two - [drumMachine](https://github.com/kdrombosky/drumMachine) and [sm-808-splice](https://github.com/vmogilev/sm-808-splice). These projects both served as inspiration for how I began implementation of some of my features, however I made sure that my project diverged in implementation from either solution early on.

### Thanks

I had a lot of fun with this coding exercise. Different aspects of the project took me to software areas I don't normally interact with (audio, user input, multi-platform, etc). Thanks for providing a fun side project to work on!