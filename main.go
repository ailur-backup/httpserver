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

func StartServer(port string, path string, address string, protocolVer string, throttleRate int64) (error, int) {
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

	var throttledFileServer http.Handler
	if throttleRate != -1 {
		throttledFileServer = ThrottleMiddleware(throttleRate)(fileServer)
	} else {
		throttledFileServer = fileServer
	}

	fmt.Println("Serving HTTP on", address, "port", port, "(http://"+address+":"+port+"/) ...")
	if protocolVer == "2.0" || protocolVer == "2" {
		httpServer = &http.Server{Addr: addressPort, Handler: throttledFileServer, TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler))}
	} else {
		httpServer = &http.Server{Addr: addressPort, Handler: throttledFileServer}
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

func ThrottleMiddleware(rate int64, burstSize ...int64) func(http.Handler) http.Handler {
	defaultBurstSize := int64(128)
	if len(burstSize) > 0 && burstSize[0] > 0 {
		defaultBurstSize = burstSize[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tw := &ThrottledResponseWriter{
				writer:    w,
				rate:      rate,
				burstSize: defaultBurstSize,
			}
			next.ServeHTTP(tw, r)
		})
	}
}

type ThrottledResponseWriter struct {
	writer    http.ResponseWriter
	rate      int64
	burstSize int64
}

func (tw *ThrottledResponseWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	var written int
	for written < len(p) {
		chunkSize := tw.burstSize
		if int64(len(p)-written) < chunkSize {
			chunkSize = int64(len(p) - written)
		}

		n, err := tw.writer.Write(p[written : written+int(chunkSize)])
		if err != nil {
			return written + n, err
		}
		written += n

		time.Sleep(time.Duration(chunkSize*8*int64(time.Second)/tw.rate) * time.Nanosecond)
	}

	return written, nil
}

func (tw *ThrottledResponseWriter) Header() http.Header {
	return tw.writer.Header()
}

func (tw *ThrottledResponseWriter) WriteHeader(statusCode int) {
	tw.writer.WriteHeader(statusCode)
}
