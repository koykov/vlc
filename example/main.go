package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/koykov/vlc"
	"log"
	"time"
)

var (
	file = flag.String("file", "", "Path to local media file")
	url  = flag.String("url", "", "URL to remote media file")
)

func main() {
	flag.Parse()

	ply, err := vlc.NewVlc([]string{})
	if err != nil {
		log.Fatal(err)
	}

	if len(*file) > 0 {
		err = ply.Play(*file)
	} else if len(*url) > 0 {
		err = ply.PlayURL(*url)
	} else {
		err = errors.New(`run command with "-h" option to see how specify input`)
	}

	if err != nil {
		log.Println(err)
	} else {
		tickDelay := time.Millisecond * 100
		tick, heartbeat := time.Tick(tickDelay), time.Tick(time.Second)
		pause0, resume0 := time.After(time.Second*10), time.After(time.Second*15)
		pause1, resume1 := time.After(time.Second*20), time.After(time.Second*25)
		finish := false
		precision := 0.01
		for {
			select {
			case <-heartbeat:
				pos, _ := ply.Position()
				fmt.Printf("pos: %f\n", pos)
			case <-tick:
				time.Sleep(tickDelay)
				pos, err := ply.Position()
				if err != nil {
					log.Fatal(err)
				}
				finish = 1.0-pos <= precision
			case <-pause0:
				fmt.Println("pause for 5 seconds")
				_ = ply.Pause()
			case <-resume0:
				fmt.Println("resume after 5 seconds")
				_ = ply.Resume()
			case <-pause1:
				fmt.Println("toggle pause for 5 seconds")
				_ = ply.TogglePause()
			case <-resume1:
				fmt.Println("toggle pause after 5 seconds")
				_ = ply.TogglePause()
			}
			if finish {
				break
			}
		}
	}

	if err := ply.Release(); err != nil {
		log.Fatal(err)
	}
}
