build:
	go mod download
	CGO_ENABLED=0 go build -ldflags "-X main.version=`git describe --tags` -X main.buildTime=`date +%FT%T%z`" -o stream2videointelligence

default: build

upgrade:
	go get -u -v
	go mod download
	go mod tidy
	go mod verify

run:
	./stream2videointelligence

clean:
	go clean
	go mod tidy
	rm -f stream2videointelligence
