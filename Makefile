.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/schedule-games schedule-games/*

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
