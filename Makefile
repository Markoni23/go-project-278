test:
	go mod tidy
	go test -v ./... -race

build:
	go build -o bin/url-shortener ./main.go

install:
	go install

lint:
	golangci-lint run ./...

air:
	air ./main.go