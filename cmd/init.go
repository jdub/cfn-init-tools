// Copyright © 2016 Jeff Waugh <jdub@bethesignal.org>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	resume  bool
	verbose bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configures a host according to stack resource metadata",
	//Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("init called")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&Stack, "configsets", "c", "default", "An optional list of configSets")

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enables verbose logging")

	if runtime.GOOS == "windows" {
		initCmd.Flags().BoolVarP(&resume, "resume", "", false, "Resume from a previous cfn-init run")
	}
}
