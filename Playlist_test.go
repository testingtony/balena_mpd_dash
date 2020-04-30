package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test the getFirstPlaylistForTime function
func TestGetFirstPlaylistForTime(t *testing.T) {
	knownWhen, expectedFirst := getFirst()

	actual := getFirstPlaylistForTime(knownWhen)

	assert.Equal(t, expectedFirst, actual)
}

func TestGetPlaylistForTime(t *testing.T) {
	knownTime, firstPlaylist := getFirst()
	tests := []struct {
		mode     modeType
		playing  string
		expected string
	}{
		// if not in playlistmode, pick the right one for the time
		{none, radio2, firstPlaylist},
		{album, radio2, firstPlaylist},

		// Playlistmode test the cycle
		{playlist, "", wantedPlaylist[0]},
		{playlist, "unknown", wantedPlaylist[0]},
		{playlist, wantedPlaylist[0], wantedPlaylist[1]},
		{playlist, wantedPlaylist[len(wantedPlaylist)-1], wantedPlaylist[0]},
	}

	for idx, test := range tests {
		playing = test.playing
		actual := getPlaylistForTime(test.mode, knownTime)
		assert.Equalf(t, test.expected, actual, "tests[%d]", idx)
	}
}
func TestGetPlaylist(t *testing.T) {
	tests := []struct {
		mode     modeType
		playing  string
		expected string
	}{
		// Playlistmode test the cycle
		{playlist, "", wantedPlaylist[0]},
		{playlist, "unknown", wantedPlaylist[0]},
		{playlist, wantedPlaylist[0], wantedPlaylist[1]},
		{playlist, wantedPlaylist[len(wantedPlaylist)-1], wantedPlaylist[0]},
	}

	for idx, test := range tests {
		playing = test.playing
		actual := getPlaylist(test.mode)
		assert.Equalf(t, test.expected, actual, "tests[%d]", idx)
	}
}

// Helpers

// a time and expected first playlist which is not in the wantedPlaylist cycle
func getFirst() (time.Time, string) {
	//Wednesday 7am = Radio X
	return time.Date(2020, time.January, 01, 7, 0, 0, 0, time.UTC), radioX
}
