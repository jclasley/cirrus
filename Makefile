build:
	go build -o .

test:
	go test server -v
	go test controller -v

all: test build