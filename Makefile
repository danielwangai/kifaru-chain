build:
	@go build -o bin/kifaru


run: build
	@./bin/kifaru

test:
	@go test -v ./...

