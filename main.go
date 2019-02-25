package main

import (
	"net/http"

	"github.com/tkivisik/books-and-albums/controllers"
)

// TODO
// * Split up the controllers
// * Health - Expose failed queries in a more interactive manner than just logs.
// * Write tests
// * Depending on requirements, expose *[]Book and *[]Album from ReturnBooksAndAlbums
//        rather than string as it is currently.

func main() {
	cultureAPIC := controllers.NewCultureAPI()

	//http.HandleFunc("/", cultureAPIC.All)
	http.Handle("/", cultureAPIC.AllMetric.TimeIt(http.HandlerFunc(cultureAPIC.All)))
	http.HandleFunc("/albums", cultureAPIC.AlbumMetric.TimeIt(cultureAPIC.Albums))
	http.HandleFunc("/books", cultureAPIC.BookMetric.TimeIt(cultureAPIC.Books))
	http.HandleFunc("/metrics", cultureAPIC.Metrics)
	http.HandleFunc("/metrics/albums", cultureAPIC.AlbumMetrics)
	http.HandleFunc("/metrics/books", cultureAPIC.BookMetrics)

	go cultureAPIC.AllMetric.Listen()
	go cultureAPIC.BookMetric.Listen()
	go cultureAPIC.AlbumMetric.Listen()

	http.ListenAndServe(":8080", nil)
}
