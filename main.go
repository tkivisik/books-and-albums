package main

import (
	"net/http"

	"github.com/tkivisik/books-and-albums/controllers"
)

// TODO
// * Split up the controllers
// * Make metrics available /metrics
// * Health - Expose failed queries in a more interactive manner than just logs.
// * Health - Keep special distinction between source API errors, my service errors.
// * Write tests
// * Depending on requirements, expose *[]Book and *[]Album from ReturnBooksAndAlbums
//        rather than string as it is currently.

func main() {
	cultureAPIC := controllers.NewCultureAPI()

	http.HandleFunc("/", cultureAPIC.All)
	http.HandleFunc("/albums", cultureAPIC.Albums)
	http.HandleFunc("/books", cultureAPIC.Books)
	//	http.HandleFunc("/metrics", cultureAPIC.Metrics) // TODO

	http.ListenAndServe(":8080", nil)
}
