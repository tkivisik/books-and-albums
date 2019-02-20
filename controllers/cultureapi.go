package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/tkivisik/books-and-albums/views"
	"google.golang.org/api/books/v1"
	"google.golang.org/api/googleapi"
)

func NewCultureAPI() *CultureAPI {
	return &CultureAPI{
		AlbumsView: views.NewView("general", "albums/albums"),
		BooksView:  views.NewView("general", "books/books"),
		AllView:    views.NewView("general", "all/all"),
	}
}

type CultureAPI struct {
	AlbumsView *views.View
	BooksView  *views.View
	AllView    *views.View
}

const help = `Solve your problem by:
 * providing a query string (e.g. localhost:8080/?query=hozier&max=7
 * visiting other endpoints:
	localhost:8080/
	localhost:8080/books
	localhost:8080/albums`

func (c *CultureAPI) All(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := ExtractQueryAndMaxResults(w, r)
	if err != nil {
		return
	}

	data := Data{
		Yield: BooksAndAlbums{
			Books:  *QueryBooks(query, maxResults),
			Albums: *QueryAlbums(query, maxResults),
		},
	}
	c.AllView.RenderHTML(w, r, data)
}

func ExtractQueryAndMaxResults(w http.ResponseWriter, r *http.Request) (string, int64, error) {
	q := r.URL.Query()
	query := q.Get("query")
	if query == "" {
		fmt.Fprintln(w, help)
		return "", 0, errors.New("no query parameter detected")
	}

	var maxResults int64
	// No max provided
	if q.Get("max") == "" {
		maxResults = 5
	} else {
		max, err := strconv.Atoi(q.Get("max"))
		if err != nil {
			fmt.Fprintf(w, "There was an error converting your provided 'max' argument: %s\n", q.Get("max"))
			return "", 0, err
		}
		maxResults = int64(max)
	}
	return query, maxResults, nil
}

func (c *CultureAPI) Books(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := ExtractQueryAndMaxResults(w, r)
	if err != nil {
		return
	}
	data := Data{Yield: *QueryBooks(query, maxResults)}
	c.BooksView.RenderHTML(w, r, data)
}

type Data struct {
	Yield interface{}
}

type BooksAndAlbums struct {
	Books  []Book
	Albums []Album
}

type Book struct {
	Title   string
	Authors []string
}

type ByTitle []Book

func (t ByTitle) Len() int           { return len(t) }
func (t ByTitle) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTitle) Less(i, j int) bool { return t[i].Title < t[j].Title }

// QueryBooks returns a *[]Books and logs request metrics like time.
// Duration is measured for retrieving and sorting the books together
// as that is the unit of work for the user of my service.
func QueryBooks(query string, maxResults int64) *[]Book {
	start := time.Now()
	books, err := getBooks(query, maxResults)
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

	log.Printf("Book request and sort time: %s\n", requestTime)
	return books
}

// AverageRequestTime returns the average of the last 10 request times
func AverageRequestTime() time.Duration {
	var sum time.Duration
	mu.Lock()
	defer mu.Unlock()

	n := int64(len(last10RequestTimes))
	if n == 0 {
		return time.Duration(0)
	}
	for _, val := range last10RequestTimes {
		sum += val
	}
	average := int64(sum) / n
	return time.Duration(average)
}

// getBooks gets the books and sorts them.
func getBooks(query string, maxResults int64) (*[]Book, error) {
	client := http.Client{}
	svc, err := books.New(&client)
	if err != nil {
		return nil, err
	}

	// Construst the VolumesListCall
	listCall := svc.Volumes.List(query)
	listCall = listCall.PrintType("books")
	listCall = listCall.MaxResults(maxResults)
	listCall = listCall.Projection("lite")
	listCall = listCall.Fields(googleapi.Field("items"))

	results, err := listCall.Do()
	if err != nil {
		return nil, err
	}

	books := []Book{}
	for _, item := range results.Items {
		book := Book{
			Title:   item.VolumeInfo.Title,
			Authors: item.VolumeInfo.Authors,
		}
		books = append(books, book)
	}
	sort.Sort(ByTitle(books))
	return &books, nil
}

func (c *CultureAPI) Albums(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := ExtractQueryAndMaxResults(w, r)
	if err != nil {
		return
	}
	data := Data{Yield: *QueryAlbums(query, maxResults)}
	c.AlbumsView.RenderHTML(w, r, data)
}

type Albumresults struct {
	Albums []Album `json:"results"`
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

	albums := Albumresults{}
	err = json.Unmarshal(b, &albums)
	if err != nil {
		log.Println(err)
	}
	sort.Sort(ByAlbumName(albums.Albums))
	return &albums.Albums, nil
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
