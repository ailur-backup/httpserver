package httpserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

func StartServer(port string, path string, address string, protocolVer string) (error, int) {
	var httpServer *http.Server
	addressPort := address + ":" + port
	fileServer := http.FileServer(http.Dir(path))
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
