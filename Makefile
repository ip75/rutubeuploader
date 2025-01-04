.PHONY:	all build clean

all: build

build:
	go build -o rutubeuploader

clean:
	rm -f rutubeuploader token.json

freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o rutubeuploader
	tar cJpvf rutubeuploader.freebsd.tar.xz rutubeuploader README.md
	rm -f rutubeuploader

mac:
	GOOS=darwin GOARCH=amd64 go build -o rutubeuploader
	tar cJpvf rutubeuploader.darwin.tar.xz rutubeuploader README.md
	rm -f rutubeuploader

linux:
	GOOS=linux GOARCH=amd64 go build -o rutubeuploader
	tar cJpvf rutubeuploader.linux.tar.xz rutubeuploader README.md
	rm -f rutubeuploader

win:
	GOOS=windows GOARCH=amd64 go build -o rutubeuploader.exe
	tar cJpvf rutubeuploader.win64.tar.xz rutubeuploader.exe README.md
	rm -f rutubeuploader.exe

dist: freebsd mac linux win
