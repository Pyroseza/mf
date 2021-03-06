package main

import (
	"flag"
	"fmt"

	"github.com/defsub/mf"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	artists *mgo.Collection
	albums  *mgo.Collection
	tracks  *mgo.Collection
	work    chan mf.Music
	done    chan string
)

func doMusic(path string, music mf.Music, err error) error {
	t, _ := music.Track()
	fmt.Printf("%s: %s / %s - %02d. %s\n", music.FileType(), music.Artist(), music.Album(), t, music.Title())
	work <- music
	return nil
}

func doArtist(music mf.Music) error {
	artist := bson.M{
		"name": music.AlbumArtist(),
		"musicbrainz": bson.M{
			"artist": music.Tags()["musicbrainz_albumartistid"],
		},
	}

	key := bson.M{
		"name": music.AlbumArtist(),
	}

	_, err := artists.Upsert(&key, &artist)

	return err
}

func doAlbum(music mf.Music) error {
	_, trackTotal := music.Track()
	_, discTotal := music.Disc()

	album := bson.M{
		"title":  music.Album(),
		"artist": music.AlbumArtist(),
		"year":   music.Year(),
		"tracks": trackTotal,
		"discs":  discTotal,
		"musicbrainz": bson.M{
			"albumartist": music.Tags()["musicbrainz_albumartistid"],
			"release":     music.Tags()["musicbrainz_albumid"],
		},
	}

	key := bson.M{
		"artist": music.AlbumArtist(),
		"title":  music.Album(),
		"year":   music.Year(),
	}

	_, err := albums.Upsert(&key, &album)

	return err
}

func doTrack(music mf.Music) error {
	trackNum, _ := music.Track()
	discNum, _ := music.Disc()

	track := bson.M{
		"artist":      music.Artist(),
		"albumartist": music.AlbumArtist(),
		"album":       music.Album(),
		"title":       music.Title(),
		"number":      trackNum,
		"disc":        discNum,
		"year":        music.Year(),
		"musicbrainz": bson.M{
			"artist":      music.Tags()["musicbrainz_artistid"],
			"albumartist": music.Tags()["musicbrainz_albumartistid"],
			"release":     music.Tags()["musicbrainz_albumid"],
		},
	}

	key := bson.M{
		"album":       music.Album(),
		"albumartist": music.AlbumArtist(),
		"year":        music.Year(),
		"number":      trackNum,
		"disc":        discNum,
	}

	_, err := tracks.Upsert(&key, &track)

	return err
}

func worker() {
	for music := range work {
		doArtist(music)
		doAlbum(music)
		doTrack(music)
	}
	done <- ""
}

func main() {
	flag.Parse()
	root := flag.Arg(0)

	url := "127.0.0.1"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// artists
	artists = session.DB("music").C("artists")
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = artists.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	// albums
	albums = session.DB("music").C("albums")
	index = mgo.Index{
		Key:        []string{"title", "artist", "year"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = albums.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	// tracks
	tracks = session.DB("music").C("tracks")
	index = mgo.Index{
		Key:        []string{"album", "albumartist", "year", "number", "disc"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = tracks.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	work = make(chan mf.Music)
	done = make(chan string)

	workers := 5

	for i := 0; i < workers; i++ {
		go worker()
	}

	if len(root) > 0 {
		mf.ScanMusic(root, doMusic)
		close(work)
	}

	for i := 0; i < workers; i++ {
		<-done
	}
}
