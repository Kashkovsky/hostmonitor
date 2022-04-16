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
	printer := core.NewPrinter()
	watcher := core.NewWatcher(&watchConfig)
	store := core.NewStore()

	go func() {
		for {
			printer.ToTable(&store)
			time.Sleep(time.Second)
		}
	}()

	watcher.Watch(store.AddOrUpdate)
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
