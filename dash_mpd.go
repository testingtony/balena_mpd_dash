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
		case amazon.Playlist:
			mpd.StopAndClear()
			mpd.AddPlaylist(getPlaylist(mode))
			mode = playlist
			mpd.Play()
		case amazon.Stop:
			mpd.StopAndClear()
			mode = none
		case amazon.Album:
			mpd.StopAndClear()
			mode = album
			mpd.AddRandomAlbum()
			mpd.Play()
		}
	}
}

const (
	radio6      string = "Radio 6 Music"
	radio2             = "Radio 2"
	radio4extra        = "Radio 4 Extra"
	radio4             = "Radio 4"
	radioX             = "Radio X"
)

var wantedPlaylist = [...]string{radio6, radio2, radio4extra, radio4}
var playing = ""

func getPlaylist(mode modeType) string {
	return getPlaylistForTime(mode, time.Now())
}

func getPlaylistForTime(mode modeType, when time.Time) string {

	if mode == playlist {
		index := -1
		for i, title := range wantedPlaylist {
			if title == playing {
				index = i
			}
		}
		index++
		if index >= len(wantedPlaylist) {
			index = 0
		}
		playing = wantedPlaylist[index]
	} else {
		playing = getFirstPlaylistForTime(when)
	}

	return playing
}

func getFirstPlaylistForTime(t time.Time) string {
	hour := t.Hour()
	switch t.Weekday() {
	case time.Sunday:
		if hour >= 10 {
			return radio2
		}
	case time.Saturday:
		if hour >= 11 && hour < 14 {
			return radio2
		}
	default:
		if hour >= 6 && hour < 10 {
			return radioX
		}
	}
	return radio6
}
