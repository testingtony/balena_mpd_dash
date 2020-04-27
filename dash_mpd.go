package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/testingtony/dash_mpd/mpd"

	"github.com/testingtony/dash_mpd/amazon"
)

type modeType int

const (
	none modeType = iota
	playlist
	album
)

func main() {

	rand.Seed(time.Now().UnixNano())

	ch, err := amazon.ConnectAndSubscribe()
	if err != nil {
		if err == amazon.ErrNoDetails {
			fmt.Println("No amazon connection details given, exiting cleanly")
			os.Exit(0)
		}
		panic(err)
	}
	mpd := mpd.NewConnection()

	mode := none
	for {
		press := <-ch
		switch press {
		case amazon.Single:
			mpd.StopAndClear()
			mpd.AddPlaylist(getPlaylist(mode))
			mode = playlist
			mpd.Play()
		case amazon.Double:
			mpd.StopAndClear()
			mode = none
		case amazon.Long:
			mpd.StopAndClear()
			mode = album
			mpd.AddRandomAlbum()
			mpd.Play()
		}
	}
}

var wantedPlaylist = [...]string{"Radio 6 Music", "Radio 2", "Radio 4 Extra", "Radio 4"}
var index int = 0

func getPlaylist(mode modeType) string {
	if mode == playlist {
		index++
		if index >= len(wantedPlaylist) {
			index = 0
		}
	} else {
		index = 0
	}

	return wantedPlaylist[index]
}
