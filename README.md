Snip helps you manage snippets in cmd.


### Some commands
```bash
# Run tests
go test ./...

# Test coverage
go test -coverprofile cover.out ./...
go tool cover -html cover.out # view as html
go tool cover -func cover.out # output coverage per functions

# Run linters
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1
golangci-lint run ./...

# build
go build -o snip .

```
