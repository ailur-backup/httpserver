# HTTPServer

This is a simple HTTP server meant to emulate python's SimpleHTTPServer module.

## Usage

```httpserver [-h] [--cgi] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [port]```

## Installing

First, have Go installed. Latest version, please. I'm talking about you debian.

Run as root

```
CDIR=$PWD
cd /tmp
git clone https://git.ailur.dev/ailur/httpserver --depth=1
cd httpserver/httpserver
make install
```

## Compiling
```
git clone https://git.ailur.dev/ailur/httpserver --depth=1
cd httpserver/httpserver
make
```
This creates the binary "httpserver"