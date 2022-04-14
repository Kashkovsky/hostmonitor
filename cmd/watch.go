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
	"time"

	"sync"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Test availability of URLs from config",
	Run:   runWatch,
}

func runWatch(cmd *cobra.Command, args []string) {
	log.Default().Println("Testing URLs from config ", watchConfig.ConfigUrl)
	resMap := sync.Map{}
	printer := core.NewPrinter()
	watcher := core.NewWatcher(&watchConfig)

	go func() {
		for {
			printer.ToTable(&resMap)
			time.Sleep(time.Second)
		}
	}()

	watcher.Watch(func(res core.TestResult) {
		if res.InProgress {
			_, ok := resMap.Load(res.Id)
			if !ok {
				resMap.Store(res.Id, res)
			}
		} else {
			resMap.Store(res.Id, res)
		}
	})
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&watchConfig.ConfigUrl, "configUrl", "c", core.ITArmyConfigURL, "Url of config containing url list")
	watchCmd.Flags().IntVarP(&watchConfig.TestInterval, "testInterval", "i", 20, "Interval in seconds between test updates")
	watchCmd.Flags().IntVarP(&watchConfig.RequestTimeout, "requestTimeout", "t", 10, "Request timeout")
	watchCmd.Flags().IntVarP(&watchConfig.UpdateInterval, "updateInterval", "u", 600, "Config update interval in seconds")
}
