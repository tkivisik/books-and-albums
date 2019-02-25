# Books and Albums

## Goal

* Test out Google Books API and iTunes API
* Output books and albums related to a search query
* Make using the two APIs independent (concurrent, resilient)
* Expose metrics / health checks

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

> **Culture API**
> **Albums**
> 1. Hozier - From Eden - EP
> 2. Hozier - Hozier (Bonus Track Version)
> 3. Hozier - Nina Cried Power - EP
> 4. Hozier - Take Me to Church - EP
> 5. Hozier - Wasteland, Baby!
> 
> **Books**
> 1. A. Hozier-Byrne - Better Love (from "The Legend of Tarzan")
> 2. Mary Soames - Clementine Churchill
> 3. \- Hozier
> 4. Perfect Papers - Keep Calm and Listen to Hozier
> 5. Henry M. Hozier - The Seven weeks' war

### GET /metrics

response (as string):
> **Culture API**  
> **Metrics**  
> Average: 124.951304ms  

### GET /metrics/books

response (as string):
> **Culture API**  
> **Book Request Metrics**  
> Average: 124.951304ms  

### GET /metrics/albums

response (as string):
> **Culture API**  
> **Album Request Metrics**  
> Average: 124.951304ms  

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
