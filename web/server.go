package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/gorilla/websocket"
)

//go:embed static
var embededFiles embed.FS

type Server struct {
	port     int
	upgrader *websocket.Upgrader
	config   *core.WatchConfig
	results  sync.Map
}

func NewServer(config *core.WatchConfig, port int) Server {
	var upgrader = websocket.Upgrader{}
	results := sync.Map{}
	return Server{port: port, upgrader: &upgrader, config: config, results: results}
}

func (s *Server) Run() {
	s.startWatch()

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
	s.sendResults(c)
}

func (s *Server) startWatch() {
	watcher := core.NewWatcher(s.config)
	go watcher.Watch(func(res core.TestResult) {
		if res.InProgress {
			_, ok := s.results.Load(res.Id)
			if !ok {
				s.results.Store(res.Id, res)
			}
		} else {
			s.results.Store(res.Id, res)
		}
	})
}

func (s *Server) sendResults(c *websocket.Conn) {
	for {
		s.results.Range(func(_ any, res interface{}) bool {
			err := c.WriteJSON(res)
			if err != nil {
				log.Println("write:", err)
				return false
			}
			return true
		})

		time.Sleep(time.Second)
	}
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
