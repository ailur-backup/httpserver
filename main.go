package httpserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Fun fact: the reference time is the specific time Jan 2 15:04:05 2006 MST, which reads:
// 01 02 03 04 05 06 -0700
// 1234567. Very clever, Google.
var timeLayout = "02/Jan/2006 15:04:05"

func StartServer(port string, path string, address string, protocolVer string) (error, int) {
	var httpServer *http.Server
	addressPort := address + ":" + port
	fileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := path + r.URL.Path
		http.ServeFile(w, r, filePath)
		var ip string
		if strings.Contains(r.RemoteAddr, ":") {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			ip = r.RemoteAddr
		}
		fmt.Println(ip + " - - [" + time.Now().Format(timeLayout) + "] \"" + r.Method + " " + r.URL.Path + " " + r.Proto + "\" " + "200" + " -")
	})

	fmt.Println("Serving HTTP on", address, "port", port, "(http://"+address+":"+port+"/) ...")
	if protocolVer == "2.0" || protocolVer == "2" {
		httpServer = &http.Server{Addr: addressPort, Handler: fileServer, TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler))}
	} else {
		httpServer = &http.Server{Addr: addressPort, Handler: fileServer}
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		if err.Error() == "permission denied" || err.Error() == "listen tcp "+addressPort+": bind: permission denied" {
			return errors.New("permission denied"), 1
		} else {
			return err, 1
		}
	}
	return nil, 0
}
