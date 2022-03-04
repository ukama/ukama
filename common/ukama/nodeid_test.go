package ukama

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ntypes := []string{"HomeNode", "homenode", "HOMENODE",
		"CompNode", "compnode", "COMPNODE",
		"AmpNode", "ampnode", "AMPNODE"}

	for _, n := range ntypes {
		nodeid := NewVirtualNodeId(n)

		_, err := ValidateNodeId(string(nodeid))
		if err != nil {
			t.Errorf("Expected Error nil; Got %s", err.Error())
		}
	}
}

func TestValidate(t *testing.T) {
	nodeid := "UK-SA2156-HNODE-A1-XXXX"

	uid, err := ValidateNodeId(string(nodeid))
	if err != nil {
		t.Errorf("Expected Error nil; Got %s", err.Error())
	}

	if !(strings.EqualFold(nodeid, string(uid))) {
		t.Errorf("Expected  %s;Got %s", strings.ToLower(nodeid), string(uid))
	}

}

func TestNegativeValidateCase1(t *testing.T) {
	nodeid := "UK-SA2156"

	_, err := ValidateNodeId(string(nodeid))
	if err == nil {
		t.Errorf("Expected Error; Got nil")
	}

}

func TestNegativeValidateCase2(t *testing.T) {
	nodeid := "UK-SA2156-CNODE-A1-XXXX"

	_, err := ValidateNodeId(string(nodeid))
	if err == nil {
		t.Errorf("Expected Error ; Got nil")
	}
}

func TestNodeType(t *testing.T) {

	ntypes := []string{"HomeNode", "homenode", "HOMENODE",
		"CompNode", "compnode", "COMPNODE",
		"AmpNode", "ampnode", "AMPNODE"}

	for _, n := range ntypes {
		nodeid := NewVirtualNodeId(n)

		res, err := ValidateNodeId(string(nodeid))
		if err != nil {
			t.Errorf("Expected Error nil; Got %s", err.Error())
		}
		assert.Equal(t, strings.ToLower(GetNodeCodeForUnits(n)), res.GetNodeType())
	}
}
