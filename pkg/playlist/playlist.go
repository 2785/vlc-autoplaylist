package playlist

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Playlist struct {
	XMLName   xml.Name `xml:"playlist"`
	Xmlns     string   `xml:"xmlns,attr"`
	Vlc       string   `xml:"vlc,attr"`
	Version   string   `xml:"version,attr"`
	Title     string   `xml:"title"`
	TrackList struct {
		Track []Track `xml:"track"`
	} `xml:"trackList"`
	Extension struct {
		Application string `xml:"application,attr"`
		Item        []struct {
			Tid string `xml:"tid,attr"`
		} `xml:"item"`
	} `xml:"extension"`
	mu sync.Mutex
}

const vlcPlaylistAppType string = "http://www.videolan.org/vlc/playlist/0"

func (pl *Playlist) Add(name string) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	escapedPath, err := getEscapedPath(name)
	if err != nil {
		return fmt.Errorf("error adding file %s: %w", name, err)
	}
	newTrack := Track{
		Location: fmt.Sprintf("file:///%s", escapedPath),
		Extension: struct {
			Application string "xml:\"application,attr\""
			ID          string "xml:\"id\""
		}{
			Application: vlcPlaylistAppType,
			ID:          strconv.Itoa(len(pl.TrackList.Track)),
		},
	}
	pl.TrackList.Track = append(pl.TrackList.Track, newTrack)
	pl.Extension.Item = append(pl.Extension.Item, struct {
		Tid string "xml:\"tid,attr\""
	}{
		Tid: strconv.Itoa(len(pl.TrackList.Track)),
	})

	return nil
}

func getEscapedPath(p string) (string, error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("cannot get full path: %w", err)
	}

	fragments := strings.Split(absPath, string(filepath.Separator))
	for i := range fragments {
		fragments[i] = url.PathEscape(fragments[i])
	}

	return strings.Join(fragments, "/"), nil
}

func (pl *Playlist) Exists(name string) bool {
	for _, v := range pl.TrackList.Track {
		if v.FileName() == name {
			return true
		}
	}
	return false
}

type Track struct {
	Location  string `xml:"location"`
	Duration  string `xml:"duration"`
	Extension struct {
		Application string `xml:"application,attr"`
		ID          string `xml:"id"`
	} `xml:"extension"`
}

func (t Track) FileName() string {
	urlEncoded := path.Base(t.Location)
	decoded, err := url.QueryUnescape(urlEncoded)
	if err != nil {
		panic("wtf")
	}

	return decoded
}
