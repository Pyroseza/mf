package mf

import (
	"os"
	"path/filepath"
	"regexp"
)

// type Media int

// const (
// 	Music Media = iota
// 	Video
// )

// type Library struct {
// 	Name string
// 	Path string
// 	Kind Media
// }

type Music interface {
	FileType() string
	Title() string
	Album() string
	Artist() string
	AlbumArtist() string
	Composer() string
	Year() int
	Genre() string
	Track() (int, int)
	Disc() (int, int)
	// Picture() *Picture
	Lyrics() string
	Tags() map[string]string
}

type Track interface {
	Track() int
	Type() string
	Codec() string
	Language() string
	Pixel() (int, int)
	Display() (int, int)
	SamplingFrequency() int
	Channels() int
}

type Video interface {
	FileType() string
	Tracks() []Track
}

type VideoFunc func(path string, video Video, err error) error

func ScanVideo(path string, scanFun VideoFunc) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		matched, err := regexp.MatchString("\\.(mkv)$", path)
		if !matched {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		video, err := parseMatroska(file)
		return scanFun(path, video, err)
	})
}

type MusicFunc func(path string, music Music, err error) error

func ScanMusic(path string, scanFn MusicFunc) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		//matched, err := regexp.MatchString("\\.(mp3|m4a|ogg)$", path)
		matched, err := regexp.MatchString("\\.(mp3)$", path)
		if !matched {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		music, err := parseMusic(file)
		return scanFn(path, music, err)
	})
}
