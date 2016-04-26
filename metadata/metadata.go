package metadata

import (
	"encoding/json"
	"errors"
)

func Parse(metadata string) (m Metadata, err error) {
	if err = json.Unmarshal([]byte(metadata), &m); err != nil {
		return
	}

	// FIXME: support configsets
	// var configs map[string]*Config ?
	if m.Init.ConfigSets != nil {
		return m, errors.New("configSets not supported yet")
	}

	return
}

type Metadata struct {
	Init *Init `json:"AWS::CloudFormation::Init"`
}

type Init struct {
	ConfigSets map[string][]interface{} `json:"configSets"`
	Config     *Config                  `json:"config"`
}

// Arranged in order of execution
type Config struct {
	Packages map[string]*Package        `json:"packages"`
	Groups   map[string]*Group          `json:"groups"`
	Users    map[string]*User           `json:"users"`
	Sources  map[string]string          `json:"sources"`
	Files    map[string]*File           `json:"files"`
	Commands map[string]*Command        `json:"commands"`
	Services map[string]*ServiceManager `json:"services"`
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
	IgnoreErrors        bool              `json:"ignoreErrors"`
	WaitAfterCompletion bool              `json:"waitAfterCompletion"`
}

type ServiceManager struct {
	SysVInit map[string]*Service `json:"sysvinit"`
	Windows  map[string]*Service `json:"windows"`
}

type Service struct {
	EnsureRunning bool                `json:"ensureRunning"`
	Enabled       bool                `json:"enabled"`
	Files         []string            `json:"files"`
	Sources       []string            `json:"sources"`
	Packages      map[string][]string `json:"packages"`
	Commands      []string            `json:"commands"`
}
