package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/2785/vlc-autoplaylist/pkg/playlist"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var plFile string
var mediaFileExtension string
var dir string
var skipConfirmation bool

func updatePlaylist(cmd *cobra.Command, args []string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(plFile)
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}

	pl := &playlist.Playlist{}
	xml.Unmarshal(b, pl)

	files, err := filepath.Glob(fmt.Sprintf("*%s", mediaFileExtension))
	if err != nil {
		return fmt.Errorf("error listing file: %w", err)
	}

	fileToAdd := []string{}

	for _, file := range files {
		if !pl.Exists(file) {
			fileToAdd = append(fileToAdd, file)
		}
	}

	if len(fileToAdd) == 0 {
		fmt.Printf("nothing to add, all present in playlist\n")
		return nil
	}

	if !skipConfirmation {
		yes, no := "yes", "no"
		fmt.Printf("Adding the following file to the playlist:\n\n  -%s\n\n", strings.Join(fileToAdd, "\n  -"))

		prompt := promptui.Select{
			Label: "Proceed",
			Items: []string{yes, no},
		}

		_, result, _ := prompt.Run()

		if result == no {
			fmt.Printf("Okay, aborting...\n")
			return nil
		}
	}

	newBytes, err := xml.Marshal(pl)
	if err != nil {
		return fmt.Errorf("Error marshaling xml: %w", err)
	}

	f, err := os.OpenFile(plFile, os.O_RDWR|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return fmt.Errorf("Error opening file for writing: %w", err)
	}
	_, err = f.Write(newBytes)

	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}

	fmt.Printf("Successfully added %v files to the playlist", len(fileToAdd))

	return nil
}
