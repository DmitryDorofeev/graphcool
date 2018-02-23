all: test

test:
	go install github.com/DmitryDorofeev/graphcool
	graphcool ./tests/models.go
	go test -v ./tests
