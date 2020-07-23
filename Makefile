build:
	cp ./.env ./bin
	env GOOS=linux go build -ldflags="-s -w" -o bin/boost-ai ./main.go

.PHONY: clean
clean:
	rm -rf ./bin/*

.PHONY: deploy
deploy: clean build
	echo 'No target yet'
