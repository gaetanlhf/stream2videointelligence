build:
	go mod download
	CGO_ENABLED=0 go build -ldflags "-X main.version=`git describe --tags` -X main.buildTime=`date +%FT%T%z`" -o cloud-video-intelligence-api-streaming

default: build

upgrade:
	go get -u -v
	go mod download
	go mod tidy
	go mod verify

run:
	./cloud-video-intelligence-api-streaming

clean:
	go clean
	go mod tidy
	rm -f cloud-video-intelligence-api-streaming
