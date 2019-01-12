# The task

Using your favorite GO framework / libraries build a service, that will accept a request with text parameter on input. It will return maximum of 5 books and maximum of 5 albums that are related to the input term. The response elements will only contain title, authors(/artists) and information whether it's a book or an album.

For albums please use the iTunes API:
https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api/#searching

For books please use Google Books API:
https://developers.google.com/books/docs/v1/reference/volumes/list

Sort the result by title alphabetically.

Make sure the software is production-ready from resilience, stability and performance point of view.

The stability of the downstream service may not be affected by the stability of the upstream services.

Results originating from one upstream service (and its stability / performance) may not affect the results originating from the other upstream service.

Your service needs to respond within a minute;

Make sure the service:
* is self-documenting
* exposes metrics on response times for upstream services
* exposes health check
* Limit of results on upstream services must be configurable per environment and preconfigured to 5.

Please document how we can run it.

Please shortly document your justification of technology / mechanism choice.

Good luck!