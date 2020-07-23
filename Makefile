build:
	cp ./.env ./bin
	env GOOS=linux go build -ldflags="-s -w" -o bin/boost-ai ./main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose --aws-profile chattabot
