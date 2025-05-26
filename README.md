# Baobab Project

A test project to use nextroute for routing clients where each client has the same
number of optional stops. Of those stops exactly one stop should be visited.

## Setup

Run the example with:

```
go run main.go -runner.input.path input.json -solve.duration 1s
```

Or build them for your platform (Mac (Apple Silicon), Linux, Windwows):

```
GOOS=darwin GOARCH=arm64 go build -o baobab-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o baobab-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o baobab-windows-amd64.exe .
```
