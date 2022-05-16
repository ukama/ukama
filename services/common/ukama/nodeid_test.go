package ukama

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ntypes := []string{"HomeNode", "homenode", "HOMENODE", "hnode", "HNODE",
		"TowerNode", "towernode", "TNODE", "tnode",
		"AmpNode", "ampnode", "AMPNODE", "ANODE", "anode"}

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

	ntypes := map[string]string{
		"HNode": NODE_ID_TYPE_HOMENODE,
		"hnode": NODE_ID_TYPE_HOMENODE,
		"HOMENODE": NODE_ID_TYPE_HOMENODE,
		"Towernode": NODE_ID_TYPE_TOWERNODE,
		"tnode": NODE_ID_TYPE_TOWERNODE,
		"TNODE": NODE_ID_TYPE_TOWERNODE,
		"ANode":  NODE_ID_TYPE_AMPNODE,
		"anode":  NODE_ID_TYPE_AMPNODE,
		"ANODE":  NODE_ID_TYPE_AMPNODE}

	for k, v := range ntypes {
		nodeid := NewVirtualNodeId(k)

		res, err := ValidateNodeId(string(nodeid))
		if assert.NoError(t, err) {
			assert.Equal(t, v, res.GetNodeType())
		}
	}
}
