test:
	go mod tidy
	go test -v ./...

build:
	go build -o bin/url-shortener ./cmd/url-shortener

install:
	go install

lint:
	golangci-lint run ./...

air:
	air ./cmd/url-shortener/main.go