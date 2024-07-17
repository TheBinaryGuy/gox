all: build tokenize

build:
	@go build -o bin/gox cmd/main.go

tokenize:
	@./bin/gox tokenize -f test.lox
