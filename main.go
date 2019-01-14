package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/tkivisik/books-and-albums/albumrequest"
	"github.com/tkivisik/books-and-albums/bookrequest"
)

// TODO
// * Remove duplication in albumrequest and bookrequest packages
// * Make album metrics available at /album/avg
// * Health - Expose failed queries in a more interactive manner than just logs.
// * Health - Keep special distinction between source API errors, my service errors.
// * Write tests
// * Depending on requirements, expose *[]Book and *[]Album from ReturnBooksAndAlbums
//        rather than string as it is currently.

// ReturnBooksAndAlbums is a high level function returning a formatted
// string including books and albums sorted alphabetically
func ReturnBooksAndAlbums(query string, maxResults int64) string {
	wg := sync.WaitGroup{}
	var books *[]bookrequest.Book

	// Make queries to Google Books API and iTunes concurrently
	wg.Add(1)
	go func() {
		books = bookrequest.QueryBooks(query, maxResults)
		wg.Done()
	}()

	var albums *[]albumrequest.Album
	wg.Add(1)
	go func() {
		albums = albumrequest.QueryAlbums(query, maxResults)
		wg.Done()
	}()

	wg.Wait()
	// TODO - use channels to communicate whichever result arrives first
	//			and print it immediate. Currently waits for both responses
	//			before writing anything.
	response := "BOOKS:\n"
	for _, book := range *books {
		response += fmt.Sprintf(" * \"%s\" by", book.Title)
		for _, author := range book.Authors {
			response += fmt.Sprintf(" %s", author)
		}
		response += "\n"
	}

	response += fmt.Sprintln("\n\nALBUMS:")
	for _, album := range *albums {
		response += fmt.Sprintf(" * \"%s\" by %s\n", album.AlbumName, album.ArtistName)
	}
	return response
}

func booksAndAlbums(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	if query == "" {
		fmt.Fprintln(w, "Solve your problem by:")
		fmt.Fprintln(w, "* providing a query string (e.g. localhost:8080/?query=hozier&max=7")
		fmt.Fprintln(w, "* visiting other endpoints:")
		fmt.Fprintln(w, "    localhost:8080/books/avg")
		fmt.Fprintln(w, "    localhost:8080/album/avg")
		return
	}

	var maxResults int64
	// No max provided
	if q.Get("max") == "" {
		maxResults = 5
	} else {
		max, err := strconv.Atoi(q.Get("max"))
		if err != nil {
			fmt.Fprintf(w, "There was an error converting your provided 'max' argument: %s\n", q.Get("max"))
			return
		}
		maxResults = int64(max)
	}
	fmt.Fprint(w, ReturnBooksAndAlbums(query, maxResults))
}

func bookAverage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Average request times for Google Books API: %s\n", bookrequest.AverageRequestTime())
}

func albumAverage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Work in Progress with album averages :)")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Solve your problem by:")
	fmt.Fprintln(w, "* providing a query string (e.g. localhost:8080/?query=hozier&max=7")
	fmt.Fprintln(w, "* visiting other endpoints:")
	fmt.Fprintln(w, "    localhost:8080/books/avg")
	fmt.Fprintln(w, "    localhost:8080/album/avg")
}

func main() {
	http.HandleFunc("/", booksAndAlbums)
	http.HandleFunc("/books/avg", bookAverage)
	http.HandleFunc("/album/avg", albumAverage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
