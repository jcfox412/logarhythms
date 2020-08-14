# LogaRhythms [![PkgGoDev](https://pkg.go.dev/badge/github.com/jcfox412/logarhythms)](https://pkg.go.dev/github.com/jcfox412/logarhythms) [![Go Report Card](https://goreportcard.com/badge/github.com/jcfox412/logarhythms)](https://goreportcard.com/report/github.com/jcfox412/logarhythms)

![rotating pattern grid with a moving cursor underneath](/docs/logarhythms_animation.gif?raw=true)

LogaRhythms is an interactive and controllable drum sequencer capable of playing 3 built-in tracks, with the ability to easily extend to more tracks.

```sh
go get -u github.com/jcfox412/logarhythms
```

## How To Run

Please make sure you've followed the Prerequisites section before trying to run LogaRhythms.

To run, simply execute the following command:

```sh
make run
```

## Prerequisite

Please make sure you have `go` installed before attempting to run.

Prerequisites for the Beep audio library come from Oto; please ensure these packages are installed before trying to run LogaRhythms.

### macOS

LogaRhythms requires `AudioToolbox.framework`, but this is automatically linked.

### Linux

libasound2-dev is required. On Ubuntu or Debian, run this command:

```sh
apt install libasound2-dev
```

In most cases this command must be run by root user or through `sudo` command.

### FreeBSD

OpenAL is required. Install openal-soft:

```sh
pkg install openal-soft
```

### OpenBSD

OpenAL is required. Install openal:

```sh
pkg_add -r openal
```

## License

[MIT](https://github.com/jcfox412/logarhythms/LICENSE)

