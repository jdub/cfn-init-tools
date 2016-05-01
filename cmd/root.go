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
	"github.com/jdub/cfn-init-tools/config"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	Config config.Config
)

var RootCmd = &cobra.Command{
	Use:   "cfn",
	Short: "A suite of host automation utilities for CloudFormation",
}

func init() {
	RootCmd.PersistentFlags().StringVar(&Config.Local, "local", "", "Local metadata JSON file")

	RootCmd.PersistentFlags().StringVarP(&Config.Stack, "stack", "s", "", "A CloudFormation stack")
	RootCmd.PersistentFlags().StringVarP(&Config.Resource, "resource", "r", "", "A CloudFormation logical resource ID")
	RootCmd.PersistentFlags().StringVar(&Config.Region, "region", "us-east-1", "The CloudFormation region")
	RootCmd.PersistentFlags().StringVarP(&Config.Url, "url", "u", "", "The CloudFormation service URL. The endpoint URL must match the region option. Use of this parameter is discouraged.")

	RootCmd.PersistentFlags().StringVarP(&Config.CredFile, "credential-file", "f", "", "OBSOLETE: Use a standard credentials file and/or AWS_PROFILE environment variable")
	RootCmd.PersistentFlags().StringVar(&Config.Role, "role", "", "OBSOLETE: IAM Role credentials will be used automatically")
	RootCmd.PersistentFlags().StringVar(&Config.AccessKey, "access-key", "", "OBSOLETE: Use a standard credentials file or AWS_ACCESS_KEY_ID environment variable")
	RootCmd.PersistentFlags().StringVar(&Config.SecretKey, "secret-key", "", "OBSOLETE: Use a standard credentials file or AWS_SECRET_ACCESS_KEY environment variable")

	RootCmd.PersistentFlags().StringVar(&Config.HttpProxy, "http-proxy", "", "A (non-SSL) HTTP proxy")
	RootCmd.PersistentFlags().StringVar(&Config.HttpsProxy, "https-proxy", "", "An HTTPS proxy")

	if runtime.GOOS == "windows" {
		Config.DataDir = os.ExpandEnv(`${SystemDrive}\cfn\cfn-init\data`)
	} else {
		Config.DataDir = "/var/lib/cfn-init/data"
	}
}
