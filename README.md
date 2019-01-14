# Books and Albums

## Usage

Spin up a server:
```bash
go run main.go
```

Visit `localhost:8080` for more info and instructions.

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
