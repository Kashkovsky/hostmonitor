/*
Copyright © 2022 Denys Kashkovskyi <dannie.k@me.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"net/http"
	"sync"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Spin up a websocket server",
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0:8080", "Server address")
	serveCmd.Flags().StringVarP(&watchConfig.ConfigUrl, "configUrl", "c", core.ITArmyConfigURL, "Url of config containing url list")
	serveCmd.Flags().IntVarP(&watchConfig.TestInterval, "testInterval", "i", 20, "Interval in seconds between test updates")
	serveCmd.Flags().IntVarP(&watchConfig.RequestTimeout, "requestTimeout", "t", 10, "Request timeout")
	serveCmd.Flags().IntVarP(&watchConfig.UpdateInterval, "updateInterval", "u", 600, "Config update interval in seconds")

}

var upgrader = websocket.Upgrader{}
var address string

func runServe(cmd *cobra.Command, args []string) {
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", serveHome)
	log.Fatal(http.ListenAndServe(address, nil))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()
	watcher := core.NewWatcher(&watchConfig)
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

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
