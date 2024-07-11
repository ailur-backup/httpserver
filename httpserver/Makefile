# Run go build
DESTDIR = /usr

all: httpserver

httpserver: main.go
	go build main.go
	mv main httpserver

clean:
	rm -f httpserver

.PHONY: all clean

install: httpserver
	cp -r httpserver $(DESTDIR)/bin

uninstall:
	rm -f $DESTDIR/bin/httpserver