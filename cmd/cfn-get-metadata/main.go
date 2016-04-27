package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/jdub/cfn-init-tools/metadata"
	"net/url"
	"os"
)

var (
	stack       string
	resource    string
	key         string
	region      string
	credfile    string
	iam_role    string
	access_key  string
	secret_key  string
	endpoint    string
	http_proxy  string
	https_proxy string

	data_dir string
)

func init() {
	flag.StringVar(&stack, "stack", "", "Name of the Stack.")
	flag.StringVar(&stack, "s", "", "Name of the Stack.")

	flag.StringVar(&resource, "resource", "", "The logical resource ID of the resource that contains the metadata.")
	flag.StringVar(&resource, "r", "", "The logical resource ID of the resource that contains the metadata.")

	flag.StringVar(&key, "key", "", "")
	flag.StringVar(&key, "k", "", "")

	flag.StringVar(&region, "region", "us-east-1", "The AWS CloudFormation regional endpoint to use.")

	flag.StringVar(&credfile, "credential-file", "", "OBSOLETE: Use a standard credentials file.")
	flag.StringVar(&credfile, "f", "", "OBSOLETE: Use a standard credentials file.")

	flag.StringVar(&iam_role, "role", "", "OBSOLETE: IAM Role credentials will be used automatically.")

	flag.StringVar(&access_key, "access-key", "", "OBSOLETE: Use a standard credentials file or AWS_ACCESS_KEY_ID environment variable.")

	flag.StringVar(&secret_key, "secret-key", "", "OBSOLETE: Use a standard credentials file or AWS_SECRET_ACCESS_KEY environment variable.")

	flag.StringVar(&endpoint, "url", "", "The AWS CloudFormation endpoint to use. Not recommended.")
	flag.StringVar(&endpoint, "u", "", "The AWS CloudFormation endpoint to use. Not recommended.")

	flag.StringVar(&http_proxy, "http-proxy", "", "An HTTP proxy (non-SSL). Use the following format: http://user:password@host:port")

	flag.StringVar(&https_proxy, "https-proxy", "", "An HTTPS proxy. Use the following format: https://user:password@host:port")
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
		return fmt.Errorf("You must specify both a stack name and logical resource id")
	}

	// FIXME: handle key
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

	json, err := metadata.ParseJson(*res.StackResourceDetail.Metadata)
	if err != nil {
		return err
	}

	fmt.Println(json)

	return nil
}
