package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/gorilla/websocket"
)

//go:embed static
var embededFiles embed.FS

type Server struct {
	port     int
	upgrader *websocket.Upgrader
	config   *core.WatchConfig
}

func NewServer(config *core.WatchConfig, port int) Server {
	var upgrader = websocket.Upgrader{}
	return Server{port: port, upgrader: &upgrader, config: config}
}

func (s *Server) Run() {
	http.HandleFunc("/ws", s.serveWs)
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	http.Handle("/", http.FileServer(getFileSystem(useOS)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil))
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()
	watcher := core.NewWatcher(s.config)
	mu := sync.Mutex{}

	watcher.Watch(func(res core.TestResult) {
		mu.Lock()
		defer mu.Unlock()
		err = c.WriteJSON(res)
		if err != nil {
			log.Println("write:", err)
		}
	})
}

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("static"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
