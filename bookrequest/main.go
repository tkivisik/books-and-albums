package bookrequest

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"google.golang.org/api/books/v1"
	"google.golang.org/api/googleapi"
)

var (
	mu                 sync.Mutex
	last10RequestTimes []time.Duration
)

type Book struct {
	Title   string
	Authors []string
}

type ByTitle []Book

func (t ByTitle) Len() int           { return len(t) }
func (t ByTitle) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTitle) Less(i, j int) bool { return t[i].Title < t[j].Title }

// QueryBooks returns a *[]Books and logs request metrics like time and success
func QueryBooks(query string, maxResults int64) (*[]Book, error) {
	start := time.Now()
	books, err := getBooks(query, maxResults)
	if err != nil {
		return nil, err
	}
	requestTime := time.Since(start)

	mu.Lock()
	last10RequestTimes = append(last10RequestTimes, requestTime)
	if len(last10RequestTimes) > 10 {
		last10RequestTimes = last10RequestTimes[1:]
	}
	mu.Unlock()

	fmt.Printf("request time: %s\n", requestTime)
	fmt.Println(books)
	return books, nil
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
