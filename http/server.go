package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	artists *mgo.Collection
	albums  *mgo.Collection
	tracks  *mgo.Collection
)

// TODO: is there a beter way?
func merge(a bson.M, b bson.M) (c bson.M) {
	c = make(bson.M)
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func artistQuery(id string) bson.M {
	if bson.IsObjectIdHex(id) {
		return bson.M{"_id": bson.ObjectIdHex(id)}
	} else {
		return bson.M{"musicbrainz.artist": id}
	}
}

func marshal(w http.ResponseWriter, data []bson.M, single bool) {
	var result []byte

	if single && len(data) >= 1 {
		result, _ = json.Marshal(data[0])
	} else {
		result, _ = json.Marshal(data)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(result)
}

func doArtists(w http.ResponseWriter, r *http.Request) {
	var data []bson.M
	var artist bson.M
	var album bson.M

	path := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")

	switch len(path) {
	case 1:
		// /artists/
		artists.Find(nil).All(&data)
		marshal(w, data, false)
	case 2:
		// /artists/id
		artists.Find(artistQuery(path[1])).Limit(1).All(&data)
		marshal(w, data, true)
	case 3:
		// /artists/id/albums
		if path[2] == "albums" {
			artists.Find(artistQuery(path[1])).One(&artist)
			albums.Find(bson.M{"artist": artist["name"]}).All(&data)
			marshal(w, data, false)
		// /artists/id/tracks
		} else if path[2] == "tracks" {
			artists.Find(artistQuery(path[1])).One(&artist)
			tracks.Find(bson.M{
				"albumartist": artist["name"],
			}).Sort("disc", "number").All(&data)
			marshal(w, data, false)
		}
	case 4:
		// /artists/id/albums/id
		if path[2] == "albums" {
			artists.Find(artistQuery(path[1])).One(&artist)
			albums.Find(merge(bson.M{"artist": artist["name"]}, albumQuery(path[3]))).Limit(1).All(&data)
			marshal(w, data, true)
		}
	case 5:
		// /artists/id/albums/id/tracks
		if path[2] == "albums" && path[4] == "tracks" {
			artists.Find(artistQuery(path[1])).One(&artist)
			albums.Find(merge(bson.M{"artist": artist["name"]}, albumQuery(path[3]))).One(&album)
			tracks.Find(bson.M{
				"album":       album["title"],
				"albumartist": album["artist"],
				"year":        album["year"]}).Sort("disc", "number").All(&data)
			marshal(w, data, false)
		}
	}
}

func albumQuery(id string) bson.M {
	if bson.IsObjectIdHex(id) {
		return bson.M{"_id": bson.ObjectIdHex(id)}
	} else {
		return bson.M{"musicbrainz.release": id}
	}
}

func doAlbums(w http.ResponseWriter, r *http.Request) {
	var data []bson.M
	var album bson.M

	path := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")

	switch len(path) {
	case 1:
		albums.Find(nil).All(&data)
		marshal(w, data, false)
	case 2:
		albums.Find(albumQuery(path[1])).Limit(1).All(&data)
		marshal(w, data, true)
	case 3:
		switch path[2] {
		case "tracks":
			albums.Find(albumQuery(path[1])).One(&album)
			tracks.Find(bson.M{
				"album":       album["title"],
				"albumartist": album["artist"],
				"year":        album["year"]}).Sort("disc", "number").All(&data)
			marshal(w, data, false)
		}
	}
}

func trackQuery(id string) bson.M {
	if bson.IsObjectIdHex(id) {
		return bson.M{"_id": bson.ObjectIdHex(id)}
	}
	return nil
}

func doTracks(w http.ResponseWriter, r *http.Request) {
	var data []bson.M

	path := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")

	switch len(path) {
	case 1:
		tracks.Find(nil).All(&data)
		marshal(w, data, false)
	case 2:
		tracks.Find(trackQuery(path[1])).Limit(1).All(&data)
		marshal(w, data, true)
	}
}

func main() {
	flag.Parse()
	// root := flag.Arg(0)

	url := "127.0.0.1"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	artists = session.DB("music").C("artists")
	albums = session.DB("music").C("albums")
	tracks = session.DB("music").C("tracks")

	http.HandleFunc("/artists/", doArtists)
	http.HandleFunc("/albums/", doAlbums)
	http.HandleFunc("/tracks/", doTracks)

	fmt.Printf("running...\n")
	http.ListenAndServe(":8080", nil)
}
