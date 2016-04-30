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
	"github.com/spf13/cobra"
)

var (
	Stack      string
	Resource   string
	Region     string
	Url        string
	HttpProxy  string
	HttpsProxy string
	CredFile   string
	Role       string
	AccessKey  string
	SecretKey  string
)

var RootCmd = &cobra.Command{
	Use:   "cfn",
	Short: "A suite of host automation utilities for CloudFormation",
}

func init() {
	//cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&Stack, "stack", "s", "", "A CloudFormation stack")
	RootCmd.PersistentFlags().StringVarP(&Resource, "resource", "r", "", "A CloudFormation logical resource ID")
	RootCmd.PersistentFlags().StringVarP(&Region, "region", "", "us-east-1", "The CloudFormation region")
	RootCmd.PersistentFlags().StringVarP(&Url, "url", "u", "", "The CloudFormation service URL. The endpoint URL must match the region option. Use of this parameter is discouraged.")

	RootCmd.PersistentFlags().StringVarP(&CredFile, "credential-file", "f", "", "OBSOLETE: Use a standard credentials file and/or AWS_PROFILE environment variable")
	RootCmd.PersistentFlags().StringVarP(&Role, "role", "", "", "OBSOLETE: IAM Role credentials will be used automatically")
	RootCmd.PersistentFlags().StringVarP(&AccessKey, "access-key", "", "", "OBSOLETE: Use a standard credentials file or AWS_ACCESS_KEY_ID environment variable")
	RootCmd.PersistentFlags().StringVarP(&SecretKey, "secret-key", "", "", "OBSOLETE: Use a standard credentials file or AWS_SECRET_ACCESS_KEY environment variable")

	RootCmd.PersistentFlags().StringVar(&HttpProxy, "http-proxy", "", "A (non-SSL) HTTP proxy")
	RootCmd.PersistentFlags().StringVar(&HttpsProxy, "https-proxy", "", "An HTTPS proxy")
}
