package core

import (
	"log"
	"net/url"
	"strings"
	"time"
)

type Watcher struct {
	config *WatchConfig
	tester Tester
	quit   chan bool
	out    chan TestResult
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
		err := w.doWatch()

		if err != nil {
			log.Fatalf("Fatal: %v", err)
			return
		}

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

		time.Sleep(time.Duration(w.config.UpdateInterval) * time.Second)
		w.quit <- true
	}
}

func (w *Watcher) doWatch() error {
	log.Default().Println("Fetching new config...")
	config, err := w.config.Update()

	if err != nil {
		log.Fatalf("Could not obtain a config: %v", err.Error())
		return err
	}

	records := strings.Split(config, "\n")
	for _, addr := range records {
		u, err := url.Parse(addr)
		if err != nil {
			log.Default().Printf("Invalid url: %v, skipping...", addr)
			continue
		}

		go w.tester.Test(u)
	}

	return nil
}
