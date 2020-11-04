# Auto VLC Playlist

This application adds media files in a folder to a playlist if they are not already present in the playlist.

## Quick Start

Say we have a folder that contains the following files:

```
- media1.mp4
- media2.mp4
- newMedia1.mp4
- newMedia2.mp4
- playlist.xspf
```

and the playlist contains `media1.mp4` and `media2.mp4`

Running `go get github.com/2785/vlc-autoplaylist` to download the app, cd into the folder, running

```
vlc-autoplaylist
```

Will add the two missing media files into the playlist.

## Usage

```
Usage:
  vlc-autoplaylist [flags]

Flags:
  -d, --directory string              The directory for the media files and the playlist (default "./")
  -h, --help                          help for vlc-autoplaylist
  -e, --media-file-extension string   The extension of media files to add to the playlist (default ".mp4")
  -p, --playlist-file string          The name of the playlist file (default "playlist.xspf")
  -y, --yes                           if need to skip confirmation
```

App supports running from a different directory, custom media extension (e.g `.avi`)
