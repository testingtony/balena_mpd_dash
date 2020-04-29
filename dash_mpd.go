package main

import (
	"fmt"
	"math/rand"
	"net"
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
	mpd, err := mpd.NewConnection()
	if err != nil {
		if dns, ok := err.(*net.DNSError); ok {
			if dns.IsNotFound {
				fmt.Println(err, "sleeping for a bit")
				time.Sleep(1 * time.Minute)
				os.Exit(1)
			}
		}
		panic(err)
	}

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
