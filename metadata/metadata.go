package metadata

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Parse(metadata string) (m Metadata, err error) {
	if err = json.Unmarshal([]byte(metadata), &m); err != nil {
		return Metadata{Init: nil}, err
	}

	if m.Init.ConfigSets != nil {
		var c ConfigSets
		if err = json.Unmarshal([]byte(metadata), &c); err != nil {
			return Metadata{Init: nil}, err
		}
		// Bring the map of configs back to the metadata return value
		m.Init.Configs = c.Configs
		// We don't want the master config when we have a map of configs
		m.Init.Config = nil
		// The map of configs should not include a configSets member
		delete(m.Init.Configs, "configSets")
	}

	return
}

type Metadata struct {
	Init *Init `json:"AWS::CloudFormation::Init"`
}

type Init struct {
	ConfigSets map[string][]interface{} `json:"configSets"`
	Config     *Config                  `json:"config"`
	Configs    map[string]*Config
}

// To fetch the map of configs when configSets != nil
type ConfigSets struct {
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
