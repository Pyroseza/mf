package mf

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/quadrifoglio/go-mkv"
)

type dimension struct {
	width  int
	height int
}

type track struct {
	number    int
	trackType string
	codec     string
	language  string
	pixel     dimension
	display   dimension
	sampling  int
	channels  int
}

func (t track) Track() int {
	return t.number
}

func (t track) Type() string {
	return t.trackType
}

func (t track) Codec() string {
	return t.codec
}

func (t track) Language() string {
	return t.language
}

func (t track) Pixel() (int, int) {
	return t.pixel.width, t.pixel.height
}

func (t track) Display() (int, int) {
	return t.display.width, t.display.height
}

func (t track) SamplingFrequency() int {
	return t.sampling
}

func (t track) Channels() int {
	return t.channels
}

type video struct {
	fileType string
	tracks []Track
}

func (v video) FileType() string {
	return v.fileType
}

func (v video) Tracks() []Track {
	return v.tracks
}

func trackType(t byte) string {
	switch t {
	case 1:
		return "VIDEO"
	case 2:
		return "AUDIO"
	case 0x11:
		return "SUBTITLE"
	}
	return ""
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func parseMatroska(r io.Reader) (video, error) {
	var err error

	tracks := make([]track, 0)
	index := -1

	doc := mkv.InitDocument(r)
	for {
		el, err := doc.ParseElement()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			break
		}
		if el.Name == "Cluster" {
			break
		}

		switch el.Name {
		case "TrackEntry":
			tracks = append(tracks, track{})
			index++
		case "TrackNumber":
			tracks[index].number = int(el.Content[0])
		case "TrackType":
			tracks[index].trackType = trackType(el.Content[0])
		case "CodecID":
			tracks[index].codec = string(el.Content)
		case "Language":
			tracks[index].language = string(el.Content)
		case "PixelWidth":
			tracks[index].pixel.width = int(binary.BigEndian.Uint16(el.Content))
		case "PixelHeight":
			tracks[index].pixel.height = int(binary.BigEndian.Uint16(el.Content))
		case "DisplayWidth":
			tracks[index].display.width = int(binary.BigEndian.Uint16(el.Content))
		case "DisplayHeight":
			tracks[index].display.height = int(binary.BigEndian.Uint16(el.Content))
		case "SamplingFrequency":
			tracks[index].sampling = int(float32frombytes(el.Content))
		case "Channels":
			tracks[index].channels = int(el.Content[0])
		}
	}

	v := video{}
	v.fileType = "MKV"
	v.tracks = make([]Track, len(tracks))
	for i, t := range tracks {
		v.tracks[i] = t
	}

	return v, err
}
