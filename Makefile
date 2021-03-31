build:
	go build -o .

test:
	go test -v

all: test build