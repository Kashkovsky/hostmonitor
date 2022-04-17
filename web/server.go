package web

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
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
	watcher  core.Watcher
}

func NewServer(config *core.WatchConfig, port int) Server {
	upgrader := websocket.Upgrader{}
	watcher := core.NewWatcher(config)
	store := core.NewStore()

	return Server{
		port:     port,
		watcher:  watcher,
		upgrader: &upgrader,
		config:   config,
		store:    store,
	}
}

func (s *Server) Run() {
	go s.watcher.Watch()
	go func() {
		for {
			select {
			case <-s.watcher.Updated:
				log.Default().Println("Clean store...")
				s.store.Clear()
			case res := <-s.watcher.Out:
				if strings.EqualFold(res.RoundId, s.watcher.RoundId.String()) {
					s.store.AddOrUpdate(res)
				}
			}
		}
	}()

	http.HandleFunc("/ws", s.serveWs)
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	http.Handle("/", http.FileServer(getFileSystem(useOS)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil))
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Creating new connection", r.RemoteAddr)
	close := make(chan bool)
	c, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	c.SetCloseHandler(func(_ int, _ string) error {
		log.Default().Println("Gracefully closing connection")
		close <- true
		return nil
	})
	go func(c *websocket.Conn) {
		c.ReadMessage()
	}(c)
	go s.sendResults(c, close)
}

func (s *Server) sendJSON(c *websocket.Conn, data interface{}, close chan bool) bool {
	err := c.WriteJSON(data)
	if err != nil {
		log.Default().Println("Error sending message", err.Error())
		close <- true
		return false
	}
	return true
}

func (s *Server) sendResults(c *websocket.Conn, close chan bool) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-close:
			c.Close()
			return
		case <-s.watcher.Updated:
			go func() {
				log.Default().Println("Sending reset message")
				s.sendJSON(c, NewResetMessage(), close)
				time.Sleep(time.Second)
			}()
		case <-tick.C:
			s.store.ForEach(func(res core.TestResult) bool {
				return s.sendJSON(c, NewResultMessage(res), close)
			})
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
