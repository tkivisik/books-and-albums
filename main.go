package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/tkivisik/books-and-albums/albumrequest"
	"github.com/tkivisik/books-and-albums/bookrequest"
)

func Do(query string, maxResults int64) {
	wg := sync.WaitGroup{}

	var books *[]bookrequest.Book
	var err error
	wg.Add(1)
	go func() {
		books, err = bookrequest.QueryBooks(query, 5)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	var albums *[]albumrequest.Album
	wg.Add(1)
	go func() {
		albums, err = albumrequest.QueryAlbums(query, 5)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("\n\nBOOKS:")
	for _, book := range *books {
		fmt.Printf(" * \"%s\" by", book.Title)
		for _, author := range book.Authors {
			fmt.Printf(" %s", author)
		}
		fmt.Println()
	}

	fmt.Println("\n\nALBUMS:")
	for _, album := range *albums {
		fmt.Printf(" * \"%s\" by %s\n", album.AlbumName, album.ArtistName)
	}

}

func main() {
	Do("Hozier", 5)
}
