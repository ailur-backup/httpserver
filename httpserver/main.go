package main

import (
	"concord.hectabit.org/HectaBit/httpserver"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func handleErr(err error, errCode int) {
	if err != nil {
		if err.Error() == "address already in use" {
			fmt.Println("OSError: [Errno 98] Address already in use")
			os.Exit(1)
		} else if err.Error() == "invalid port" {
			fmt.Println("socket.gaierror: [Errno -8] Servname not supported for ai_socktype")
			os.Exit(1)
		} else if err.Error() == "permission denied" {
			fmt.Println("PermissionError: [Errno 13] Permission denied")
			os.Exit(1)
		} else {
			fmt.Println("Unknown error:", err)
			os.Exit(errCode)
		}
	} else {
		os.Exit(errCode)
	}
}

func main() {
	if len(os.Args) < 2 {
		err, errCode := httpserver.StartServer("8000", "./", "0.0.0.0", "1.0", -1)
		handleErr(err, errCode)
	} else {
		args := os.Args[1:]
		protocolVer, path, address, port := "1.0", "./", "0.0.0.0", "8000"
		throttleRate := int64(-1)
		for num, arg := range args {
			if strings.Contains(arg, "-") {
				if strings.Contains(arg, "-p") || strings.Contains(arg, "--protocol") {
					if len(args) > num+1 {
						protocolVer = args[num+1]
						args = append(args[:num+1], args[num+2:]...)
					} else {
						fmt.Println("usage: httpserver [-h] [--cgi] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\nhttpserver: error: argument -p/--protocol: expected one argument")
						os.Exit(2)
					}
				} else if strings.Contains(arg, "-d") || strings.Contains(arg, "--directory") {
					if len(args) > num+1 {
						path = args[num+1]
						args = append(args[:num+1], args[num+2:]...)
					} else {
						fmt.Println("usage: httpserver [-h] [--cgi] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\nhttpserver: error: argument -d/--directory: expected one argument")
						os.Exit(2)
					}
				} else if strings.Contains(arg, "-b") || strings.Contains(arg, "--bind") {
					if len(args) > num+1 {
						address = args[num+1]
						args = append(args[:num+1], args[num+2:]...)
					} else {
						fmt.Println("usage: httpserver [-h] [--cgi] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\nhttpserver: error: argument -b/--bind: expected one argument")
						os.Exit(2)
					}
				} else if strings.Contains(arg, "-t") || strings.Contains(arg, "--throttle") {
					if len(args) > num+1 {
						var err error
						throttleRate, err = strconv.ParseInt(args[num+1], 10, 64)
						if err != nil {
							fmt.Println("usage: httpserver [-h] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\nhttpserver: error: argument -t/--throttle: invalid int value: '" + args[num+1] + "'")
							os.Exit(2)
						}
					} else {
						fmt.Println("usage: httpserver [-h] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\nhttpserver: error: argument -t/--throttle: expected one argument")
						os.Exit(2)
					}
				} else if strings.Contains(arg, "-h") || strings.Contains(arg, "--help") {
					fmt.Println("usage: httpserver [-h] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]\n\npositional arguments:\n  port                  bind to this port (default: 8000)\n\noptions:\n  -h, --help            show this help message and exit\n  -b ADDRESS, --bind ADDRESS\n                        bind to this address (default: all interfaces)\n  -d DIRECTORY, --directory DIRECTORY\n                        serve this directory (default: current directory)\n  -p VERSION, --protocol VERSION\n                        conform to this HTTP version (default: HTTP/1.0)")
					os.Exit(0)
				} else {
					fmt.Println("usage: httpserver [-h] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]")
					fmt.Println("httpserver: error: unrecognized arguments: " + arg)
					os.Exit(2)
				}
			} else {
				if len(args) >= num+1 {
					if args[num] == arg {
						_, err := strconv.Atoi(arg)
						if err != nil {
							fmt.Println("usage: httpserver [-h] [-b ADDRESS] [-d DIRECTORY] [-p VERSION] [-t RATE] [port]")
							fmt.Println("httpserver: error: argument port: invalid int value: '" + arg + "'")
							os.Exit(2)
						} else {
							port = arg
						}
					}
				}
			}
		}
		err, errCode := httpserver.StartServer(port, path, address, protocolVer, throttleRate)
		handleErr(err, errCode)
	}
}
