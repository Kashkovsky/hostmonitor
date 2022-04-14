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
	"os"

	"github.com/Kashkovsky/hostmonitor/core"
	"github.com/spf13/cobra"
)

var watchConfig = core.WatchConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hostmonitor",
	Short: "A simple utility to monitor host availability by given config.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&watchConfig.ConfigUrl, "configUrl", "c", core.ITArmyConfigURL, "Url of config containing url list")
	rootCmd.PersistentFlags().IntVarP(&watchConfig.TestInterval, "testInterval", "i", 20, "Interval in seconds between test updates")
	rootCmd.PersistentFlags().IntVarP(&watchConfig.RequestTimeout, "requestTimeout", "t", 10, "Request timeout")
	rootCmd.PersistentFlags().IntVarP(&watchConfig.UpdateInterval, "updateInterval", "u", 600, "Config update interval in seconds")
}
