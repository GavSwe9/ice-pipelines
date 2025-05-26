.PHONY: build clean deploy

build:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/populate-game-queue/bootstrap populate-game-queue/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/process-game/bootstrap process-game/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/process-player-season-totals/bootstrap process-player-season-totals/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/process-player/bootstrap process-player/*

zip:
	zip -j build/lambda/populate-game-queue.zip build/lambda/populate-game-queue/bootstrap
	zip -j build/lambda/process-game.zip build/lambda/process-game/bootstrap
	zip -j build/lambda/process-player-season-totals.zip build/lambda/process-player-season-totals/bootstrap
	zip -j build/lambda/process-player.zip build/lambda/process-player/bootstrap

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
