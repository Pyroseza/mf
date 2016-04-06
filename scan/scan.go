package main

import (
	"flag"
	"fmt"

	"github.com/defsub/mf"
)

func doMusic(path string, music mf.Music, err error) error {
	t, _ := music.Track()
	fmt.Printf("%s: %s / %s - %02d %s\n", music.FileType(), music.Artist(), music.Album(), t, music.Title())
	fmt.Printf("%s\n", music.Tags()["musicbrainz_albumid"])
	return nil
}

func doVideo(path string, video mf.Video, err error) error {
	fmt.Printf("%+v\n", path)
	fmt.Printf("%+v\n", video.FileType())
	for _, t := range video.Tracks() {
		fmt.Printf(" %d %s %s %s\n", t.Track(), t.Type(), t.Codec(), t.Language())
	}
	return nil
}

func main() {
	flag.Parse()
	root := flag.Arg(0)

	//mf.ScanVideo(root, doVideo)

	if len(root) > 0 {
		// lib := mf.Library{"My Music", root, mf.Music}
		mf.ScanMusic(root, doMusic)
	}
}
