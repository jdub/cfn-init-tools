package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/davecgh/go-spew/spew"
	"github.com/jdub/cfn-init-tools/metadata"
	"net/url"
	"os"
	"runtime"
)

var (
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

	if stack == "" || resource == "" {
		return errors.New("You must specify both a stack name and logical resource id")
	}

	// FIXME: should we provide a workaround for aws-sdk-go's AWS_PROFILE vs. standard AWS_DEFAULT_PROFILE?
	// FIXME: handle http/https_proxy

	config := aws.NewConfig()

	config.Region = aws.String(region)

	if endpoint != "" {
		if u, err := url.Parse(endpoint); err != nil {
			return err
		} else if u.Scheme == "" {
			return errors.New("invalid endpoint url")
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

	metadata, err := metadata.Parse(*res.StackResourceDetail.Metadata)
	if err != nil {
		return err
	}

	spew.Dump(metadata)
	return nil
}
