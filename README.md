# Books and Albums

## Project Structure

> .  
> ├── albumrequest  
> │   └── main.go  
> ├── bookrequest  
> │   └── main.go  
> ├── main.go  
> └── README.md  

## Usage

Spin up a server:
```bash
go run main.go
```

Visit `localhost:8080` for more info and instructions.

## API Contract

Current responses are given as string text. Future improvements might include responses as JSON.

### GET /

```
parameters {
    query - what are you looking for in books and albums
    max - max number of results (default 5)
}
```

e.g. `GET /?query=hoziert&max=5`

response (as string):
> BOOKS:
>  * "Better Love (from "The Legend of Tarzan")" by A. Hozier-Byrne
>  * "Hozier" by
>  * "The Seven Weeks' War" by H.M. Hozier
>  * "War in the East" by Quintin Barry
>  * "Winston and Clementine" by Mary Soames
> 
> 
> ALBUMS:
>  * "From Eden - EP" by Hozier
>  * "Hozier" by Hozier
>  * "Hozier (Bonus Track Version)" by Hozier
>  * "Nina Cried Power - EP" by Hozier
>  * "Take Me to Church - EP" by Hozier

### GET /books/avg

response (as string):
> Average request times for Google Books API: `<average based on last 10 requests>`ms

### GET /album/avg

Work in progress

---


## Testing

Currently, no unit tests have been written. Race conditions have been explored using:

```bash
go run -race main.go
```

## Cleaning

circleci-lint was used as a static analysis tool. True positives were dealt with.

```bash
circleci-lint run
```

## Further Development

See also the inline notes in code.

* Add tests
* Code review
* Integrate with CICD

* Further nail down requirements
* Based on requirements, prioritize new performance improvements / new features / usability etc.

* Give API response as JSON/XML
