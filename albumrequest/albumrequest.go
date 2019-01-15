package albumrequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

type Albums struct {
	Results []Album `json:"results"`
}

type Album struct {
	ArtistName string `json:"artistName"`
	AlbumName  string `json:"collectionName"`
}

type ByAlbumName []Album

func (a ByAlbumName) Len() int           { return len(a) }
func (a ByAlbumName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlbumName) Less(i, j int) bool { return a[i].AlbumName < a[j].AlbumName }

var (
	mu                 sync.Mutex
	last10RequestTimes []time.Duration
)

// getAlbums retrieves and unmarshals the album data from itunes
func getAlbums(query string, maxResults int64) (*[]Album, error) {
	safeQuery := url.QueryEscape(query)
	fullURL := fmt.Sprintf("https://itunes.apple.com/search?term=%s&media=music&limit=%d&entity=album", safeQuery, maxResults)
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	albums := Albums{}
	err = json.Unmarshal(b, &albums)
	if err != nil {
		log.Println(err)
	}
	sort.Sort(ByAlbumName(albums.Results))
	return &albums.Results, nil
}

// QueryAlbums returns a *[]Album and logs request metrics like time
func QueryAlbums(query string, maxResults int64) *[]Album {
	start := time.Now()
	albums, err := getAlbums(query, maxResults)
	if err != nil {
		log.Println(err)
		return nil
	}
	requestTime := time.Since(start)

	mu.Lock()
	last10RequestTimes = append(last10RequestTimes, requestTime)
	if len(last10RequestTimes) > 10 {
		last10RequestTimes = last10RequestTimes[1:]
	}
	mu.Unlock()

	log.Printf("Album request and sort time: %s\n", requestTime)
	return albums
}
