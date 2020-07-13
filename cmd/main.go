package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	// HTTP Listen port
	listen = flag.String("listen", "0.0.0.0:80", "HTTP Listen configuration")
)

// ServerMUX dumb implementation
type MyMux struct {
	lastReqID uint
	mtx       sync.RWMutex
}

func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allocate new request ID
	m.mtx.RLock()
	m.lastReqID++
	reqID := m.lastReqID
	m.mtx.RUnlock()

	fmt.Printf("[%d][%s][%s][%s]\n", reqID, r.Method, r.Host, r.URL)
	for k, v := range r.Header {
		fmt.Printf("[%d]# %s:%s\n", reqID, k, v)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[%d]! Error reading body\n")
		return
	}
	fmt.Printf("[%d]> (%d) [%s]\n\n", reqID, len(body), body)
}

//
// Simple HTTP Debug Server
//

func main() {
	flag.Parse()
	fmt.Printf("HTTP Debug Server - listening on %s\n", *listen)

	mux := MyMux{}
	err := http.ListenAndServe(*listen, &mux)
	fmt.Println("HTTP Listen error: ", err)

}
