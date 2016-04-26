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
)

var (
	stack       string
	resource    string
	region      string
	credfile    string
	configsets  string
	endpoint    string
	http_proxy  string
	https_proxy string
	verbose     bool
)

func init() {
	flag.StringVar(&stack, "stack", "", "Name of the Stack.")
	flag.StringVar(&stack, "s", "", "Name of the Stack.")

	flag.StringVar(&resource, "resource", "", "The logical resource ID of the resource that contains the metadata.")
	flag.StringVar(&resource, "r", "", "The logical resource ID of the resource that contains the metadata.")

	flag.StringVar(&region, "region", "us-east-1", "The AWS CloudFormation regional endpoint to use.")

	flag.StringVar(&credfile, "credential-file", "", "A file that contains both a secret access key and an access key.")
	flag.StringVar(&credfile, "f", "", "A file that contains both a secret access key and an access key.")

	flag.StringVar(&configsets, "configsets", "default", "A comma-separated list of configsets to run (in order).")
	flag.StringVar(&configsets, "c", "default", "A comma-separated list of configsets to run (in order).")

	flag.StringVar(&endpoint, "url", "", "The AWS CloudFormation endpoint to use.")
	flag.StringVar(&endpoint, "u", "", "The AWS CloudFormation endpoint to use.")

	flag.StringVar(&http_proxy, "http-proxy", "", "An HTTP proxy (non-SSL). Use the following format: http://user:password@host:port")
	flag.StringVar(&https_proxy, "https-proxy", "", "An HTTPS proxy. Use the following format: https://user:password@host:port")

	flag.BoolVar(&verbose, "v", false, "Verbose output. This is useful for debugging cases where cfn-init is failing to initialize.")
}

func main() {
	//fmt.Println("Os.Args:", os.Args)
	flag.Parse()

	fmt.Println("Variables")
	fmt.Println("       stack: ", stack)
	fmt.Println("    resource: ", resource)
	fmt.Println("      region: ", region)
	fmt.Println("    credfile: ", credfile)
	fmt.Println("  configsets: ", configsets)
	fmt.Println("         url: ", endpoint)
	fmt.Println("  http_proxy: ", http_proxy)
	fmt.Println(" https_proxy: ", https_proxy)
	fmt.Println("     verbose: ", verbose)

	config := aws.NewConfig()

	config.Region = aws.String(region)

	if endpoint != "" {
		if u, err := url.Parse(endpoint); err != nil {
			fmt.Println(err.Error())
		} else if u.Scheme == "" {
			fmt.Println(errors.New("invalid endpoint url"))
			return
		} else {
			config.Endpoint = aws.String(u.String())
		}
	}

	svc := cloudformation.New(session.New(), config)

	params := &cloudformation.DescribeStackResourceInput{LogicalResourceId: aws.String(resource), StackName: aws.String(stack)}

	res, err := svc.DescribeStackResource(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	metadata, err := metadata.Parse(*res.StackResourceDetail.Metadata)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	spew.Dump(metadata)
}
