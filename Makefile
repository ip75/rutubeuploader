.PHONY:	all build clean

all: build

build:
	go build -o rutubeuploader

clean:
	rm -f rutubeuploader token.json