package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	store    core.Store
}

func NewServer(config *core.WatchConfig, port int) Server {
	var upgrader = websocket.Upgrader{}
	store := core.NewStore()
	return Server{port: port, upgrader: &upgrader, config: config, store: store}
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
	close := make(chan bool)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	c.SetCloseHandler(func(_ int, _ string) error {
		close <- true
		return nil
	})
	s.sendResults(c, close)
}

func (s *Server) startWatch() {
	watcher := core.NewWatcher(s.config)
	go watcher.Watch(s.store.AddOrUpdate)
}

func (s *Server) sendResults(c *websocket.Conn, close chan bool) {
Loop:
	for {
		select {
		case <-close:
			break Loop
		default:
			s.store.ForEach(func(res core.TestResult) bool {
				err := c.WriteJSON(res)
				if err != nil {
					log.Println("Error sending a message to client", err.Error())
					close <- true
					return false
				}
				return true
			})

			time.Sleep(time.Second)
		}
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
