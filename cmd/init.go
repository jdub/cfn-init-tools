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
	"github.com/davecgh/go-spew/spew"
	"github.com/jdub/cfn-init-tools/metadata"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	configSets []string
	resume     bool
	verbose    bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configures a host according to stack resource metadata",
	//Long:  `...`,
	RunE: cfnInit,
}

func init() {
	RootCmd.AddCommand(initCmd)

	s := initCmd.Flags().StringP("configsets", "c", "default", "An optional list of configSets")
	configSets = strings.Split(*s, ",")

	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enables verbose logging")

	if runtime.GOOS == "windows" {
		initCmd.Flags().BoolVar(&resume, "resume", false, "Resume from a previous cfn-init run")
	}
}

func cfnInit(cmd *cobra.Command, args []string) error {
	if Config.Local == "" && (Config.Stack == "" || Config.Resource == "") {
		return fmt.Errorf("You must pass --local, or --stack and --resource")
	}

	raw, err := metadata.Fetch(Config)
	if err != nil {
		return err
	}

	meta, err := metadata.Parse(raw)
	if err != nil {
		return err
	}

	// Prepare the data directory for logging and whatnot
	if err := os.MkdirAll(Config.DataDir, 0644); err != nil {
		//fmt.Fprintf(os.Stderr, "Error: Could not create data directory: %v\n", Config.DataDir)
		return err
	}

	// Write fetched metadata to file
	json, err := metadata.ParseJson(raw, "")
	if err != nil {
		return err
	}

	name := filepath.Join(Config.DataDir, "metadata.json")
	if err := ioutil.WriteFile(name, []byte(json), os.FileMode(0644)); err != nil {
		//fmt.Fprintf(os.Stderr, "Error: Failed to write %v\n", name)
		return err
	}

	spew.Dump(meta)

	return nil
}
