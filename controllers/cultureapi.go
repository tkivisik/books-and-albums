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

	"github.com/tkivisik/books-and-albums/metrics"
	"github.com/tkivisik/books-and-albums/views"
	"google.golang.org/api/books/v1"
	"google.golang.org/api/googleapi"
)

// NewCultureAPI returns a pointer to a CultureAPI Controller
func NewCultureAPI() *CultureAPI {
	return &CultureAPI{
		AlbumsView:       views.NewView("general", "albums/albums"),
		BooksView:        views.NewView("general", "books/books"),
		AllView:          views.NewView("general", "all/all"),
		MetricsView:      views.NewView("general", "metrics/metrics"),
		AlbumMetricsView: views.NewView("general", "metrics/albums"),
		BookMetricsView:  views.NewView("general", "metrics/books"),
		AlbumMetric:      metrics.NewMetric(10),
		BookMetric:       metrics.NewMetric(10),
		AllMetric:        metrics.NewMetric(10),
	}
}

// CultureAPI is a controller struct
type CultureAPI struct {
	AlbumsView       *views.View
	BooksView        *views.View
	AllView          *views.View
	MetricsView      *views.View
	AlbumMetricsView *views.View
	BookMetricsView  *views.View
	AlbumMetric      *metrics.Metric
	BookMetric       *metrics.Metric
	AllMetric        *metrics.Metric
}

var ErrNoQuery = errors.New("no query parameter detected")

// All handles a request to books and albums APIs
func (c *CultureAPI) All(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := extractQueryAndMaxResults(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	var b []Book
	var a []Album
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		books, err := GetBooks(query, maxResults)
		if err != nil {
			log.Println(err)
		}
		b = *books
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		albums, err := GetAlbums(query, maxResults)
		if err != nil {
			log.Println(err)
		}
		a = *albums
		wg.Done()
	}()
	wg.Wait()

	data := Data{
		Yield: BooksAndAlbums{
			Books:  b,
			Albums: a,
		},
	}
	c.AllView.RenderHTML(w, r, data)
}

// extractQueryAndMaxResults parses the request parameters
func extractQueryAndMaxResults(w http.ResponseWriter, r *http.Request) (string, int64, error) {
	q := r.URL.Query()
	query := q.Get("query")
	if query == "" {
		fmt.Fprintln(w, help)
		return "", 0, ErrNoQuery
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

// Books handles Book queries
// GET /books?query=something
func (c *CultureAPI) Books(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := extractQueryAndMaxResults(w, r)
	if err != nil {
		return
	}
	books, err := GetBooks(query, maxResults)
	if err != nil {
		log.Println(err)
	}
	data := Data{Yield: *books}
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

// GetBooks gets the books and sorts them.
func GetBooks(query string, maxResults int64) (*[]Book, error) {
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

// Albums handles the album requests
// GET /albums?query=something
func (c *CultureAPI) Albums(w http.ResponseWriter, r *http.Request) {
	query, maxResults, err := extractQueryAndMaxResults(w, r)
	if err != nil {
		return
	}
	albums, err := GetAlbums(query, maxResults)
	if err != nil {
		log.Println(err)
	}

	data := Data{Yield: *albums}
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

// GetAlbums retrieves and unmarshals the album data from itunes
func GetAlbums(query string, maxResults int64) (*[]Album, error) {
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

func (c *CultureAPI) Metrics(w http.ResponseWriter, r *http.Request) {
	c.AllMetric.GetAvg <- true
	mean := <-c.AllMetric.Avg
	data := Data{Yield: mean}
	c.MetricsView.RenderHTML(w, r, data)
}

func (c *CultureAPI) AlbumMetrics(w http.ResponseWriter, r *http.Request) {
	c.AlbumMetric.GetAvg <- true
	mean := <-c.AlbumMetric.Avg
	data := Data{Yield: mean}
	c.AlbumMetricsView.RenderHTML(w, r, data)
}

func (c *CultureAPI) BookMetrics(w http.ResponseWriter, r *http.Request) {
	c.BookMetric.GetAvg <- true
	mean := <-c.BookMetric.Avg
	data := Data{Yield: mean}
	c.BookMetricsView.RenderHTML(w, r, data)
}

const help = `Solve your problem by:
 * providing a query string (e.g. localhost:8080/?query=hozier&max=7
 * visiting other endpoints:
	localhost:8080/
	localhost:8080/books
	localhost:8080/albums
	localhost:8080/metrics
	localhost:8080/metrics/albums
	localhost:8080/metrics/books
	`
