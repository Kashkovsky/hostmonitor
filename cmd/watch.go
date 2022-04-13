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
	"net/url"
	"strings"
	"time"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/spf13/cobra"
)

type WatchConfig struct {
	configUrl      string
	testInterval   int64
	requestTimeout int
}

var watchConfig = WatchConfig{}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Test availability of URLs from config",
	Run:   runWatch,
}

func runWatch(cmd *cobra.Command, args []string) {
	log.Default().Println("Testing URLs from config ", watchConfig.configUrl)
	printer := core.NewPrinter()
	for {
		res, err := doWatch()
		if err != nil {
			log.Fatalf("Fatal: %v", err)
			return
		}
		printer.ToTable(&res)
		time.Sleep(time.Duration(watchConfig.testInterval) * time.Second)
	}
}

func doWatch() ([]core.TestResult, error) {
	config, err := core.GetStringFromURL(watchConfig.configUrl)
	if err != nil {
		log.Fatalf("Could not obtain a config: %v", err.Error())
		return []core.TestResult{}, err
	}

	records := strings.Split(config, "\n")
	outC := make(chan core.TestResult, 50)
	for _, addr := range records {
		u, err := url.Parse(addr)
		if err != nil {
			log.Default().Printf("Invalid url: %v, skipping...", addr)
			continue
		}

		go core.Test(u, watchConfig.requestTimeout, outC)
	}

	results := []core.TestResult{}

	for {
		results = append(results, <-outC)
		if len(results) == len(records) {
			return results, nil
		}
	}
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&watchConfig.configUrl, "configUrl", "c", core.ITArmyConfigURL, "Url of config containing url list")
	watchCmd.Flags().Int64VarP(&watchConfig.testInterval, "testInterval", "i", 10, "Interval in seconds between test updates")
	watchCmd.Flags().IntVarP(&watchConfig.requestTimeout, "requestTimeout", "t", 5, "Request timeout")
}
