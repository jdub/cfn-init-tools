package metadata

import (
	"testing"
)

func TestInit(t *testing.T) {
	json := `
{
    "AWS::CloudFormation::Init": {}
}
`
	if _, err := Parse(json); err != nil {
		t.Error(err)
	}
	return
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
	} else {
		if m.Init.Config == nil {
			t.Errorf("config")
		}
	}
	return
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
                        "enabled": 1
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

		if en := m.Init.Config.Services.SysVInit["nginx"].Enabled; en != true {
			t.Errorf("%+v not interpreted as false", en)
		}
	}
	return
}
