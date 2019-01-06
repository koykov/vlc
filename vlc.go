package vlc

// #cgo LDFLAGS: -lvlc
// #include <vlc/vlc.h>
// #include <stdlib.h>
import "C"
import (
	"errors"
	"os"
	"unsafe"
)

type Vlc struct {
	instance *C.libvlc_instance_t
	player   *C.libvlc_media_player_t
	media    *C.libvlc_media_t
}

// The constructor.
func NewVlc(args []string) (*Vlc, error) {
	// Convert arguments.
	argc := C.int(len(args))
	argv := make([]*C.char, 0)
	for _, arg := range args {
		argv = append(argv, C.CString(arg))
	}
	defer func() {
		for i := range argv {
			C.free(unsafe.Pointer(argv[i]))
		}
	}()

	// Make instance.
	vlc := Vlc{
		instance: C.libvlc_new(C.int(argc), *(***C.char)(unsafe.Pointer(&argv))),
	}
	if err := vlc.getLastErr(); err != nil {
		return nil, err
	}

	// Make player.
	vlc.player = C.libvlc_media_player_new(vlc.instance)
	if err := vlc.getLastErr(); err != nil {
		return nil, err
	}

	return &vlc, nil
}

// Get latest error from libvlc.
func (vlc *Vlc) getLastErr() error {
	if err := C.libvlc_errmsg(); err != nil {
		return errors.New(C.GoString(err))
	}
	return nil
}

// Release VLC resources.
func (vlc *Vlc) Release() error {
	C.libvlc_media_player_release(vlc.player)
	vlc.player = nil

	C.libvlc_release(vlc.instance)
	vlc.instance = nil

	C.libvlc_media_release(vlc.media)
	vlc.media = nil

	return vlc.getLastErr()
}

// Play VLC media.
func (vlc *Vlc)playMedia(media *C.libvlc_media_t) error {
	C.libvlc_media_player_set_media(vlc.player, vlc.media)
	if err := vlc.getLastErr(); err != nil {
		return err
	}
	if C.libvlc_media_player_play(vlc.player) < 0 {
		return vlc.getLastErr()
	}
	return nil
}

// Play local media file.
func (vlc *Vlc) Play(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return errors.New("not found: " + filepath)
	}

	if vlc.media != nil {
		C.libvlc_media_release(vlc.media)
		vlc.media = nil
	}

	filepathPtr := C.CString(filepath)
	defer C.free(unsafe.Pointer(filepathPtr))

	if vlc.media = C.libvlc_media_new_path(vlc.instance, filepathPtr); vlc.media != nil {
		return vlc.playMedia(vlc.media)
	}

	return vlc.getLastErr()
}

// Play remote file.
func (vlc *Vlc) PlayURL(url string) error {
	if vlc.media != nil {
		C.libvlc_media_release(vlc.media)
		vlc.media = nil
	}

	urlPtr := C.CString(url)
	defer C.free(unsafe.Pointer(urlPtr))

	if vlc.media = C.libvlc_media_new_location(vlc.instance, urlPtr); vlc.media != nil {
		return vlc.playMedia(vlc.media)
	}

	return vlc.getLastErr()
}

// Returns position of played media.
func (vlc *Vlc) Position() (float64, error) {
	return float64(C.libvlc_media_player_get_position(vlc.player)), vlc.getLastErr()
}

// Pause playing.
func (vlc * Vlc) Pause() error {
	C.libvlc_media_player_set_pause(vlc.player, C.int(1))
	return vlc.getLastErr()
}

// Resume playing.
func (vlc * Vlc) Resume() error {
	C.libvlc_media_player_set_pause(vlc.player, C.int(0))
	return vlc.getLastErr()
}

// Toggle pause.
func (vlc *Vlc) TogglePause() error {
	C.libvlc_media_player_pause(vlc.player)
	return vlc.getLastErr()
}

// Stop playing.
func (vlc *Vlc) Stop() error {
	C.libvlc_media_player_stop(vlc.player)
	return vlc.getLastErr()
}
