package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/zserge/lorca"
)

//go:embed www
var fs embed.FS

func main() {
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("", "", 1000, 850)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go func() {
		log.Print("Serving www on " + ln.Addr().String())
		fs := http.FileServer(http.FS(fs))
		http.Serve(ln, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Add("Cache-Control", "no-cache")
			if strings.HasSuffix(req.URL.Path, ".wasm.gz") {
				resp.Header().Set("Content-Encoding", "gzip")
				resp.Header().Set("Content-Type", "application/wasm")
			} else if strings.HasSuffix(req.URL.Path, ".wasm") {
				resp.Header().Set("Content-Type", "application/wasm")
			}
			fs.ServeHTTP(resp, req)
		}))
	}()

	ui.Load(fmt.Sprintf("http://%s/www", ln.Addr()))

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
