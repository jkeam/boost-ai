all: test

.PHONY: test
test:
	go test

.PHONY: coverage
coverage:
  go test  -cover