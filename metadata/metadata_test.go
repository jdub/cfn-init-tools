package metadata

import (
	//"github.com/davecgh/go-spew/spew"
	"errors"
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
			t.Error(errors.New("config"))
		}
	}
	return
}
