package mf

import (
	"io"

	"github.com/dhowden/tag"
	"github.com/dhowden/tag/mbz"
)

type music struct {
	metadata tag.Metadata
	mbz mbz.Info
}

func (m music) FileType() string {
	return string(m.metadata.FileType())
}

func (m music) Title() string {
	return m.metadata.Title()
}

func (m music) Album() string {
	return m.metadata.Album()
}

func (m music) Artist() string {
	return m.metadata.Artist()
}

func (m music) AlbumArtist() string {
	return m.metadata.AlbumArtist()
}

func (m music) Composer() string {
	return m.metadata.Composer()
}

func (m music) Year() int {
	return m.metadata.Year()
}

func (m music) Genre() string {
	return m.metadata.Genre()
}

func (m music) Track() (int, int) {
	return m.metadata.Track()
}

func (m music) Disc() (int, int) {
	return m.metadata.Disc()
}

func (m music) Lyrics() string {
	return m.metadata.Lyrics()
}

func (m music) Tags() map[string]string {
	return m.mbz
}

func parseMusic(r io.ReadSeeker) (music, error) {
	m, err := tag.ReadFrom(r)
	if err != nil {
		return music{}, err
	}
	i := mbz.Extract(m)
	return music{m, i}, nil
}
