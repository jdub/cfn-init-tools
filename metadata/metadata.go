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

package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/jdub/cfn-init-tools/config"
	"io/ioutil"
	"net/url"
	"strings"
)

func Fetch(conf config.Config) (metadata string, err error) {
	if conf.Local != "" {
		if b, err := ioutil.ReadFile(conf.Local); err != nil {
			return "", err
		} else {
			return string(b), nil
		}
	}

	endpoint := ""
	if conf.Url != "" {
		if u, err := url.Parse(conf.Url); err != nil {
			return "", err
		} else if u.Scheme == "" {
			return "", fmt.Errorf("invalid endpoint url: %v", conf.Url)
		} else {
			endpoint = u.String()
		}
	}

	svc := cloudformation.New(session.New(), &aws.Config{
		Region:   aws.String(conf.Region),
		Endpoint: aws.String(endpoint),
	})

	params := &cloudformation.DescribeStackResourceInput{
		LogicalResourceId: aws.String(conf.Resource),
		StackName:         aws.String(conf.Stack),
	}

	res, err := svc.DescribeStackResource(params)
	if err != nil {
		return "", err
	}

	return *res.StackResourceDetail.Metadata, nil
}

func Parse(metadata string) (m Metadata, err error) {
	bytes := []byte(metadata)
	if err = json.Unmarshal(bytes, &m); err != nil {
		return
	}

	if m.Init == nil {
		err = fmt.Errorf("Could not find 'AWS::CloudFormation::Init' key in metadata")
		return
	}

	var c Configs
	if err = json.Unmarshal(bytes, &c); err != nil {
		return
	}
	// Bring the map of configs back to the metadata return value
	m.Init.Configs = c.Configs
	// The map of configs should not include a configSets member
	delete(m.Init.Configs, "configSets")

	// FIXME: we should validate the configsets
	return
}

func Json(metadata string, key string) (j string, err error) {
	bytes := []byte(metadata)

	var d map[string]interface{}
	if err := json.Unmarshal(bytes, &d); err != nil {
		return "", err
	}

	// FIXME: only return value at dotted object path, e.g. config.commands

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type Metadata struct {
	Authentication map[string]*Authentication `json:"AWS::CloudFormation::Authentication"`
	Init           *Init                      `json:"AWS::CloudFormation::Init"`
}

type Authentication struct {
	AccessKeyId string   `json:"accessKeyId"`
	Buckets     []string `json:"buckets"`
	Password    string   `json:"password"`
	SecretKey   string   `json:"secretKey"`
	Type        string   `json:"type"`
	Uris        []string `json:"uris"`
	Username    string   `json:"username"`
	RoleName    string   `json:"roleName"`
}

// Skips Configs which will be picked up on the second run
type Init struct {
	ConfigSets map[string][]interface{} `json:"configSets"`
	Configs    map[string]*Config       `json:"-"`
}

// To fetch the map of configs when configSets != nil
type Configs struct {
	Configs map[string]*Config `json:"AWS::CloudFormation::Init"`
}

// Arranged in order of execution
type Config struct {
	Packages map[string]*Package `json:"packages"`
	Groups   map[string]*Group   `json:"groups"`
	Users    map[string]*User    `json:"users"`
	Sources  map[string]string   `json:"sources"`
	Files    map[string]*File    `json:"files"`
	Commands map[string]*Command `json:"commands"`
	Services *ServiceManager     `json:"services"`
}

type Package struct {
	Msi      map[string]string   `json:"msi"`
	Rpm      map[string]string   `json:"rpm"`
	Yum      map[string][]string `json:"yum"`
	Apt      map[string][]string `json:"apt"`
	Python   map[string][]string `json:"python"`
	RubyGems map[string][]string `json:"rubygems"`
}

type Group struct {
	Gid string `json:"gid"`
}

type User struct {
	Uid     string   `json:"uid"`
	Groups  []string `json:"groups"`
	HomeDir string   `json:"homeDir"`
}

type File struct {
	Content        string                     `json:"content"`
	Source         string                     `json:"source"`
	Encoding       string                     `json:"encoding"`
	Group          string                     `json:"group"`
	Owner          string                     `json:"owner"`
	Mode           string                     `json:"mode"`
	Authentication string                     `json:"authentication"`
	Context        map[string]json.RawMessage `json:"context"`
}

type Command struct {
	Command             string            `json:"command"`
	Env                 map[string]string `json:"env"`
	Cwd                 string            `json:"cwd"`
	Test                string            `json:"test"`
	IgnoreErrors        JavaScriptBoolean `json:"ignoreErrors"`
	WaitAfterCompletion JavaScriptBoolean `json:"waitAfterCompletion"`
}

type ServiceManager struct {
	SysVInit map[string]*Service `json:"sysvinit"`
	Windows  map[string]*Service `json:"windows"`
}

type Service struct {
	EnsureRunning JavaScriptBoolean   `json:"ensureRunning"`
	Enabled       JavaScriptBoolean   `json:"enabled"`
	Files         []string            `json:"files"`
	Sources       []string            `json:"sources"`
	Packages      map[string][]string `json:"packages"`
	Commands      []string            `json:"commands"`
}

type JavaScriptBoolean bool

func (bit *JavaScriptBoolean) UnmarshalJSON(data []byte) error {
	s := strings.ToLower(strings.Trim(string(data), `"`))
	if s == "1" || s == "true" {
		*bit = true
	} else if s == "0" || s == "false" {
		*bit = false
	} else {
		return fmt.Errorf("Boolean unmarshal error: invalid input %s", s)
	}
	return nil
}
