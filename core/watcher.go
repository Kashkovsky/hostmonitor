package core

import (
	"log"
	"net/url"
	"strings"
	"time"
)

type Watcher struct {
	config  *WatchConfig
	rawUrls string
	tester  Tester
	quit    chan bool
	out     chan TestResult
}

func NewWatcher(config *WatchConfig) Watcher {
	outC := make(chan TestResult, 50)
	quit := make(chan bool)
	tester := NewTester(config, outC, quit)
	return Watcher{
		config: config,
		tester: tester,
		quit:   quit,
		out:    outC,
	}
}

func (w *Watcher) Watch(f func(TestResult)) {
	for {
		err := w.update()
		if err != nil {
			log.Fatalln("Cannot proceed without config, terminating")
			return
		}
		err = w.doWatch(f)

		if err != nil {
			log.Fatalf("Fatal: %v", err)
			return
		}

		time.Sleep(time.Duration(w.config.UpdateInterval) * time.Second)
		w.quit <- true
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
	}

	return nil
}

func (w *Watcher) doWatch(f func(TestResult)) error {
	records := strings.Split(w.rawUrls, "\n")
	for _, addr := range records {
		u, err := url.Parse(addr)
		if err != nil {
			log.Default().Printf("Invalid url: %v, skipping...", addr)
			continue
		}

		go w.tester.Test(u)
		go func() {
			for {
				select {
				case <-w.quit:
					return
				case rec := <-w.out:
					f(rec)
				}
			}
		}()
	}

	return nil
}
