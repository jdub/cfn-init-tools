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
	"github.com/jdub/cfn-init-tools/metadata"
	"github.com/spf13/cobra"
)

var (
	key string
)

// getMetadataCmd represents the get-metadata command
var getMetadataCmd = &cobra.Command{
	Use:   "get-metadata",
	Short: "Fetch the metadata associated with a specified stack resource",
	//Long:  `...`,
	RunE: getMetadata,
}

func init() {
	RootCmd.AddCommand(getMetadataCmd)

	getMetadataCmd.Flags().StringVarP(&key, "key", "k", "", "Retrieve the value at <key> in the Metadata object; must be in dotted object notation (parent.child.leaf)")
}

func getMetadata(cmd *cobra.Command, args []string) error {
	meta, err := metadata.Fetch(Config)
	if err != nil {
		return err
	}

	json, err := metadata.ParseJson(meta, key)
	if err != nil {
		return err
	}

	fmt.Println(json)

	return nil
}
