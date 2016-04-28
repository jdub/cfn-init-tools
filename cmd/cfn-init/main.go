package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/jdub/cfn-init-tools/metadata"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	local       string
	stack       string
	resource    string
	region      string
	credfile    string
	iam_role    string
	access_key  string
	secret_key  string
	configsets  string
	endpoint    string
	http_proxy  string
	https_proxy string
	verbose     bool
	resume      bool

	data_dir string
)

func init() {
	flag.StringVar(&local, "local", "", "A local metadata JSON file")

	flag.StringVar(&stack, "stack", "", "Name of the Stack.")
	flag.StringVar(&stack, "s", "", "Name of the Stack.")

	flag.StringVar(&resource, "resource", "", "The logical resource ID of the resource that contains the metadata.")
	flag.StringVar(&resource, "r", "", "The logical resource ID of the resource that contains the metadata.")

	flag.StringVar(&region, "region", "us-east-1", "The AWS CloudFormation regional endpoint to use.")

	flag.StringVar(&credfile, "credential-file", "", "OBSOLETE: Use a standard credentials file.")
	flag.StringVar(&credfile, "f", "", "OBSOLETE: Use a standard credentials file.")

	flag.StringVar(&iam_role, "role", "", "OBSOLETE: IAM Role credentials will be used automatically.")

	flag.StringVar(&access_key, "access-key", "", "OBSOLETE: Use a standard credentials file or AWS_ACCESS_KEY_ID environment variable.")

	flag.StringVar(&secret_key, "secret-key", "", "OBSOLETE: Use a standard credentials file or AWS_SECRET_ACCESS_KEY environment variable.")

	flag.StringVar(&configsets, "configsets", "default", "A comma-separated list of configsets to run (in order).")
	flag.StringVar(&configsets, "c", "default", "A comma-separated list of configsets to run (in order).")

	flag.StringVar(&endpoint, "url", "", "The AWS CloudFormation endpoint to use. Not recommended.")
	flag.StringVar(&endpoint, "u", "", "The AWS CloudFormation endpoint to use. Not recommended.")

	flag.StringVar(&http_proxy, "http-proxy", "", "An HTTP proxy (non-SSL). Use the following format: http://user:password@host:port")

	flag.StringVar(&https_proxy, "https-proxy", "", "An HTTPS proxy. Use the following format: https://user:password@host:port")

	flag.BoolVar(&verbose, "v", false, "Verbose output. This is useful for debugging cases where cfn-init is failing to initialize.")

	if runtime.GOOS == "windows" {
		flag.BoolVar(&resume, "resume", false, "Resume from a previous cfn-init run.")
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	var input string

	if local != "" {
		if b, err := ioutil.ReadFile(local); err != nil {
			return err
		} else {
			input = string(b)
		}

	} else {
		if stack == "" || resource == "" {
			return fmt.Errorf("You must specify both a stack name and logical resource id")
		}

		// FIXME: should we provide a workaround for aws-sdk-go's AWS_PROFILE vs. standard AWS_DEFAULT_PROFILE?
		// FIXME: handle http/https_proxy

		config := aws.NewConfig()

		config.Region = aws.String(region)

		if endpoint != "" {
			if u, err := url.Parse(endpoint); err != nil {
				return err
			} else if u.Scheme == "" {
				return fmt.Errorf("invalid endpoint url: %v", endpoint)
			} else {
				config.Endpoint = aws.String(u.String())
			}
		}

		svc := cloudformation.New(session.New(), config)

		params := &cloudformation.DescribeStackResourceInput{LogicalResourceId: aws.String(resource), StackName: aws.String(stack)}

		res, err := svc.DescribeStackResource(params)
		if err != nil {
			return err
		}

		input = *res.StackResourceDetail.Metadata
	}

	meta, err := metadata.Parse(input)
	if err != nil {
		return err
	}

	// Prepare the data_dir
	if runtime.GOOS == "windows" {
		data_dir = os.ExpandEnv(`${SystemDrive}\cfn\cfn-init\data`)
	} else {
		data_dir = "/var/lib/cfn-init/data"
	}

	if err := os.MkdirAll(data_dir, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create data directory: %v\n", data_dir)
	} else {
		// Write fetched metadata to file
		json, err := metadata.ParseJson(input)
		if err != nil {
			return err
		}

		name := filepath.Join(data_dir, "metadata.json")
		if err := ioutil.WriteFile(name, []byte(json), os.FileMode(0644)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write %v\n", name)
			return err
		}
	}

	for name, attr := range meta.Init.Configs["config"].Files {
		if !filepath.IsAbs(name) {
			fmt.Fprintf(os.Stderr, "Error: File specified with non-absolute path: %v\n", name)
			continue
		}

		parent := filepath.Dir(name)
		if err := os.MkdirAll(parent, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not create parent directory: %v\n", parent)
			continue
		}

		mode, _ := FileModeFromString(attr.Mode)
		symlink := mode&os.ModeSymlink == os.ModeSymlink

		if symlink {
			if attr.Content == "" {
				fmt.Fprintf(os.Stderr, "Error: Symbolic link specified without a destination")
				continue
			}

			i, err := os.Lstat(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Could not get info on file: %v", name)
				continue
			}

			if i.Mode()&os.ModeSymlink != os.ModeSymlink {
				// unlink and write symlink
			}

			l, err := os.Readlink(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Could not determine destination for symlink: %v", name)
				continue
			}

			if l != attr.Content {
				// now compare then do a write temp / rename
			}
		}

		var d []byte
		if attr.Encoding == "base64" {
			d, err = base64.StdEncoding.DecodeString(attr.Content)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Could not decode base64 content for file: %v", name)
				continue
			}
		} else {
			d = []byte(attr.Content)
		}

		if err := ioutil.WriteFile(name, d, mode); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write %v\n", name)
			continue
		}

		fmt.Printf("Wrote: %v\n", name)
	}

	//spew.Dump(meta)
	return nil
}

func FileModeFromString(perm string) (mode os.FileMode, err error) {
	if perm == "" {
		return os.FileMode(0), fmt.Errorf("empty perm")
	}

	m, err := strconv.ParseInt(perm, 8, 32)
	if err != nil {
		return os.FileMode(0), fmt.Errorf("unable to parse perm")
	}

	return os.FileMode(m), nil
}
