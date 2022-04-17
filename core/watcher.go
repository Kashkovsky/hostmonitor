package core

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Watcher struct {
	Out     chan TestResult
	Updated chan bool
	config  *WatchConfig
	rawUrls string
	tester  Tester
	quit    chan bool
	RoundId *uuid.UUID
}

func NewWatcher(config *WatchConfig) Watcher {
	outC := make(chan TestResult, 50)
	quit := make(chan bool)
	updated := make(chan bool)
	tester := NewTester(config, outC)
	return Watcher{
		Updated: updated,
		config:  config,
		tester:  tester,
		quit:    quit,
		Out:     outC,
	}
}

func (w *Watcher) Watch() {
	for {
		quit := make(chan bool)
		err := w.update()
		id := uuid.New()
		w.RoundId = &id

		if err != nil {
			log.Fatalln("Cannot proceed without config, terminating")
			return
		}
		err = w.doWatch(quit)

		if err != nil {
			log.Fatalf("Fatal: %v", err)
			return
		}

		time.Sleep(time.Duration(w.config.UpdateInterval) * time.Second)
		close(quit)
	}
}

func (w *Watcher) update() error {
	log.Default().Println("Fetching new config...")
	config, err := w.config.UpdateURLs()

	if err != nil {
		log.Fatalf("Could not obtain a config: %v", err.Error())
		if w.rawUrls == "" {
			return err
		}
	}

	if strings.EqualFold(w.rawUrls, config) {
		log.Default().Println("URLs didn't change")
	} else {
		w.rawUrls = config
		log.Default().Println("New URLs have been applied")
		close(w.Updated)
		w.Updated = make(chan bool)
	}

	return nil
}

func (w *Watcher) doWatch(quit <-chan bool) error {
	records := ParceUrls(w.rawUrls)
	for _, addr := range records {
		go w.tester.Test(*w.RoundId, addr, quit)
	}

	return nil
}
