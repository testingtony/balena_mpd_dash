package mpd

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/tatsushid/go-fastping"
	"github.com/vincent-petithory/mpdclient"
)

/*
Mpd connection mode
*/
type Mpd struct {
	playlistMode bool
	albumMode    bool
	mpd          mpdclient.MPDClient
}

/*
NewConnection gets a new MPD connection
*/
func NewConnection() (*Mpd, error) {
	host, ok := os.LookupEnv("MPDHOST")
	if !ok {
		host = "localhost"
	}

	hadToWait, err := waitForMPD(host)
	if err != nil {
		return nil, err
	}
	// wait for MPD to fully start if it wasn't there at the start
	if hadToWait {
		time.Sleep(10 * time.Second)
	}

	conn, err := mpdclient.Connect(host, 6600)
	if err != nil {
		return nil, err
	}
	return &Mpd{
		playlistMode: false,
		albumMode:    false,
		mpd:          *conn,
	}, nil
}

/*
Play starts playing whatever is in the queue
*/
func (m *Mpd) Play() {
	m.mpd.Cmd("play")
}

/*
Stop stops the playing
*/
func (m *Mpd) Stop() {
	m.mpd.Cmd("stop")
}

/*
Clear empties the queue
*/
func (m *Mpd) Clear() {
	m.mpd.Cmd("clear")
}

/*
StopAndClear stops and clears
*/
func (m *Mpd) StopAndClear() {
	m.Stop()
	m.Clear()
}

/*
AddRandomAlbum chooses a reandom album from the database and adds it to the queue
*/
func (m *Mpd) AddRandomAlbum() error {
	result := m.mpd.Cmd("list album")
	if result.Err != nil {
		return result.Err
	}

	albums := make([]string, 0, len(result.Data))
	var albumRegexp = regexp.MustCompile(`Album:\s*(.+)`)
	for _, text := range result.Data {
		match := albumRegexp.FindStringSubmatch(text)
		if match == nil {
			albums = append(albums, "")
		} else {
			albums = append(albums, match[1])
		}
	}

	album := albums[rand.Intn(len(albums))]
	if err := m.addAlbum(album); err != nil {
		return err
	}

	return nil
}

func (m *Mpd) addAlbum(albumName string) error {

	result := m.mpd.Cmd(fmt.Sprintf(`find "(album == \"%s\" )"`, albumName))
	if result.Err != nil {
		return result.Err
	}
	fmt.Println("Adding Album", albumName)
	var responseRegexp = regexp.MustCompile(`file:\s*(.+)`)
	for _, result := range result.Data {
		match := responseRegexp.FindStringSubmatch(result)
		if match != nil {
			m.mpd.Cmd(fmt.Sprintf(`add "%s"`, match[1]))
		}
	}
	return nil
}

/*
AddPlaylist adds the named saved playlist to the queue
*/
func (m *Mpd) AddPlaylist(playlist string) error {

	fmt.Println("Adding playlist", playlist)
	response := m.mpd.Cmd(fmt.Sprintf(`load "%s"`, playlist))
	return response.Err
}

func waitForMPD(hostname string) (bool, error) {
	p := fastping.NewPinger()
	p.MaxRTT = 5 * time.Second
	ra, err := net.ResolveIPAddr("ip4:icmp", hostname)
	if err != nil {
		return false, err
	}

	ping := make(chan bool)

	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ping <- true
	}
	p.OnIdle = func() {
		ping <- false
	}
	p.RunLoop()

	hadToWait := false
	ticker := time.NewTicker(5 * time.Second)
	for alive := false; !alive; {
		select {
		case x := <-p.Done():
			fmt.Println(x)
			if err := p.Err(); err != nil {
				return false, err
			}
		case alive = <-ping:
			if !alive {
				hadToWait = true
			}
		case <-ticker.C:
			fmt.Println("No response from", hostname, "waiting")
		}
	}
	p.Stop()
	return hadToWait, nil
}
