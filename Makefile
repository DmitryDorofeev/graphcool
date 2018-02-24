all: test

generate:
	go install github.com/DmitryDorofeev/graphcool
	graphcool ./tests/models.go

test: generate
	go test -v ./tests
