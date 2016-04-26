package metadata

import (
	"testing"
)

func TestParseJson(t *testing.T) {
	json := `{"pants": "shirt"}`
	pretty := `{
  "pants": "shirt"
}`
	if j, err := ParseJson(json); err != nil {
		t.Error(err)
	} else if j != pretty {
		t.Errorf("%v != %v", j, pretty)
	}
}

func TestEmpty(t *testing.T) {
	json := `{}`
	// No, Mr. Bond, I expect you to die!
	if _, err := Parse(json); err != nil {
		return
	}
}

func TestInit(t *testing.T) {
	json := `
{
    "AWS::CloudFormation::Init": {}
}
`
	if _, err := Parse(json); err != nil {
		t.Error(err)
	}
}

func TestConfig(t *testing.T) {
	json := `
{
    "AWS::CloudFormation::Init": {
        "config": {}
    }
}
`
	if m, err := Parse(json); err != nil {
		t.Error(err)
	} else if m.Init.Config == nil {
		t.Errorf("Init.Config not unmarshalled correctly")
	}
}

func TestUnmarshalTruthyJSON(t *testing.T) {
	json := `
{
    "AWS::CloudFormation::Init": {
        "config": {
            "commands": {
                "ps afx": {
                    "command": "ps afx",
                    "ignoreErrors": "true",
                    "waitAfterCompletion": false
                }
            },
            "services": {
                "sysvinit": {
                    "nginx": {
                        "enabled": 1,
                        "ensureRunning": "0"
                    }
                }
            }
        }
    }
}
`
	if m, err := Parse(json); err != nil {
		t.Error(err)
	} else {
		if ie := m.Init.Config.Commands["ps afx"].IgnoreErrors; ie != true {
			t.Errorf("%+v not interpreted as true", ie)
		}

		if wac := m.Init.Config.Commands["ps afx"].WaitAfterCompletion; wac != false {
			t.Errorf("%+v not interpreted as false", wac)
		}

		if e := m.Init.Config.Services.SysVInit["nginx"].Enabled; e != true {
			t.Errorf("%+v not interpreted as true", e)
		}

		if er := m.Init.Config.Services.SysVInit["nginx"].EnsureRunning; er != false {
			t.Errorf("%+v not interpreted as false", er)
		}
	}
}

func TestConfigSets(t *testing.T) {
	json := `
{
    "AWS::CloudFormation::Init": {
        "configSets": {
            "ascending": [ "1", "2" ],
            "descending": [ "2", "1" ],
            "test": [ "test" ],
            "default": [ { "ConfigSet": "ascending" } ]
        },
        "config": {},
        "1": {},
        "2": {},
        "test": {
            "services": {
                "sysvinit": {
                    "nginx": {}
                }
            }
        }
    }
}
`
	if m, err := Parse(json); err != nil {
		t.Error(err)
	} else {
		if m.Init.Config != nil {
			t.Errorf("Init.Config should be nil when processing configSets")
		}
		if _, ok := m.Init.Configs["configSets"]; ok {
			t.Errorf(`Init.Configs["configSets"] should be nil when processing configSets`)
		}
		if _, ok := m.Init.Configs["test"].Services.SysVInit["nginx"]; !ok {
			t.Errorf(`Init.Configs["test"].Services.SysVInit["nginx"] not unmarshalled correctly`)
		}
		// FIXME: test configSets themselves
	}
}
