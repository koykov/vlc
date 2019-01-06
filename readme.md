VLC
===

Simple CGO wrapper over libvlc. Provides Go wrappers for libvlc functionality.

## Installation

This package requires installed libvlc header files. Just install it using the following command:
```bash
sudo dnf install vlc-devel
```
Then use `go get github.com/koykov/vlc` to install the package.

## Usage

```go
// See `vlc --help` or `cvlc --help` for list of available options and/or arguments. 
args := []string{"--option0", "--option1"}

player := vlc.NewVlc([]string{})
player.Play("path/to/media/file.mp3") // or PlayURL("url to remote file or stream")

time.Sleep(trackDuration)

player.Release()
```

See example directory to more detailed examples of usage.
