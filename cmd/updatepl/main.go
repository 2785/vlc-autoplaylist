package updatepl

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/2785/auto-playlist-vlc/pkg/playlist"
)

func main() {
	b, err := ioutil.ReadFile("playlist.xspf")
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		os.Exit(1)
	}

	pl := &playlist.Playlist{}
	xml.Unmarshal(b, pl)

	files, err := filepath.Glob("*.mp4")
	if err != nil {
		fmt.Printf("error listing file: %s", err)
		os.Exit(1)
	}

	for _, file := range files {
		if !pl.Exists(file) {
			err := pl.Add(file)
			if err != nil {
				fmt.Printf("error adding file '%s' to playlist: %s", file, err)
				os.Exit(1)
			}
		}
	}

	newBytes, err := xml.Marshal(pl)
	if err != nil {
		fmt.Printf("Error marshaling xml: %s", err)
		os.Exit(1)
	}

	f, err := os.OpenFile("playlist.xspf", os.O_RDWR|os.O_TRUNC, os.ModePerm)

	if err != nil {
		fmt.Printf("Error opening file for writing: %s", err)
		os.Exit(1)
	}
	_, err = f.Write(newBytes)

	if err != nil {
		fmt.Printf("Error writing file: %s", err)
		os.Exit(1)
	}

}
