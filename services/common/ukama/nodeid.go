package ukama

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	CODE_IDX     = 10
	NodeIDLength = 23
	OEMCODE      = "UK"
	MFGCODE      = "SA"
	DELIMITER    = "-"
	HWVERSION    = "M0"
	UNITVERSION  = "V0"
)

const (
	NODE_ID_TYPE_HOMENODE  = "HNODE"
	NODE_ID_TYPE_TOWERNODE = "TNODE"
	NODE_ID_TYPE_COMPNODE  = "CNODE"
	NODE_ID_TYPE_AMPNODE   = "ANODE"
	NODE_ID_TYPE_UNDEFINED = "UNDEFINED"
)

const (
	MODULE_ID_TYPE_COMP      = "COMV1"
	MODULE_ID_TYPE_TRX       = "TRX"
	MODULE_ID_TYPE_CTRL      = "CTRL"
	MODULE_ID_TYPE_FE        = "FE"
	MODULE_ID_TYPE_UNDEFINED = "UNDEFINED"
)

const (
	node_id_type_component_home      = "hnode"
	node_id_type_component_tower     = "tnode"
	node_id_type_component_amplifier = "anode"
)

type NodeID string
type ModuleID string

func (n NodeID) String() string {
	return string(n)
}

func (n NodeID) StringLowercase() string {
	return strings.ToLower(n.String())
}

func (m ModuleID) String() string {
	return string(m)
}

func (m ModuleID) StringLowercase() string {
	return strings.ToLower(m.String())
}

func (n NodeID) GetNodeType() string {
	t := n.String()[CODE_IDX : CODE_IDX+strings.IndexRune(n.String()[CODE_IDX:], '-')]
	switch strings.ToLower(t) {
	case node_id_type_component_home:
		return NODE_ID_TYPE_HOMENODE

	case node_id_type_component_tower:
		return NODE_ID_TYPE_COMPNODE

	case node_id_type_component_amplifier:
		return NODE_ID_TYPE_AMPNODE
	default:
		return NODE_ID_TYPE_UNDEFINED
	}
}

func getRandCode(t time.Time) string {
	rand.Seed(time.Now().UnixNano())
	min := 0x0000
	max := 0xFFFF
	val := rand.Intn(max-min+1) + min
	hexcode := fmt.Sprintf("%04X", val)
	return hexcode
}

/* Get HW Code */
func GetNodeCodeForUnits(ntype string) string {
	var code string
	switch ntype {
	case NODE_ID_TYPE_HOMENODE, "HomeNode", "homenode", "hnode", "HNode":
		code = "HNODE"
	case NODE_ID_TYPE_TOWERNODE, "TowerNode", "towernode", "tnode", "TNode":
		code = "TNODE"
	case NODE_ID_TYPE_AMPNODE, "AmpNode", "ampnode":
		code = "ANODE"
	default:
		code = "XXXXX"
	}
	return code
}

/* Get HW Code */
func GetModuleCodeForUnits(mtype string) string {
	var code string
	switch mtype {
	case MODULE_ID_TYPE_COMP, "ComV1", "com", "COM":
		code = "COM"
	case MODULE_ID_TYPE_TRX, "Trx", "trx":
		code = "TRX"
	case MODULE_ID_TYPE_CTRL, "Ctrl", "ctrl":
		code = "CTRL"
	case MODULE_ID_TYPE_FE, "Fe", "fe":
		code = "RFFE"
	default:
		code = "XXXXX"
	}
	return code
}

/* Generate new node id for virtual node */
func NewVirtualNodeId(ntype string) NodeID {
	t := time.Now()
	year, week := t.ISOWeek()
	yearstr := strconv.Itoa(year)
	yearcode := yearstr[len(yearstr)-2:]
	weekstr := fmt.Sprintf("%02d", week)
	code := GetNodeCodeForUnits(ntype)

	/*2+1+6+1+5+1+2+1+4*/
	/* UK-SA2154-HNODE-A1-XXXX*/
	uuid := OEMCODE + DELIMITER + MFGCODE + yearcode + weekstr + DELIMITER + code + DELIMITER + UNITVERSION + DELIMITER + getRandCode(t)

	log.Infof("UUID: New NodeID for %s is %s and length is %d", ntype, uuid, len(uuid))

	/* RFC 1123 lowercase id and tags*/
	lid := strings.ToLower(uuid)

	return NodeID(lid)
}

/* Generate new module id for virtual module */
func NewVirtualModuleId(mtype string) ModuleID {
	t := time.Now()
	year, week := t.ISOWeek()
	yearstr := strconv.Itoa(year)
	yearcode := yearstr[len(yearstr)-2:]
	weekstr := fmt.Sprintf("%02d", week)
	code := GetModuleCodeForUnits(mtype)

	/*2+1+6+1+5+1+2+1+4*/
	/* UK-SA2154-HNODE-A1-XXXX*/
	uuid := OEMCODE + DELIMITER + MFGCODE + yearcode + weekstr + DELIMITER + code + DELIMITER + HWVERSION + DELIMITER + getRandCode(t)

	log.Infof("UUID: New ModuleID for %s is %s and length is %d", mtype, uuid, len(uuid))

	/* RFC 1123 lowercase id and tags*/
	lid := strings.ToLower(uuid)

	return ModuleID(lid)
}

func NewVirtualComId() ModuleID {
	return NewVirtualModuleId(MODULE_ID_TYPE_COMP)
}

func NewVirtualTRXId() ModuleID {
	return NewVirtualModuleId(MODULE_ID_TYPE_TRX)
}

func NewVirtualCtrlId() ModuleID {
	return NewVirtualModuleId(MODULE_ID_TYPE_CTRL)
}

func NewVirtualFEId() ModuleID {
	return NewVirtualModuleId(MODULE_ID_TYPE_FE)
}

// Generate new node id for home node
func NewVirtualHomeNodeId() NodeID {
	return NewVirtualNodeId(NODE_ID_TYPE_HOMENODE)
}

// Generate new node id for amplifier node
func NewVirtualAmplifierNodeId() NodeID {
	return NewVirtualNodeId(NODE_ID_TYPE_AMPNODE)
}

// Generate new node id for tower node
func NewVirtualTowerNodeId() NodeID {
	return NewVirtualNodeId(NODE_ID_TYPE_TOWERNODE)
}

func ValidateNodeId(id string) (NodeID, error) {

	/* TODO :: ADD more validation once we finalized this format */
	if len(id) != NodeIDLength {
		err := errors.New("invalid length")
		return "", err
	}

	/* Check for HW codes */
	codes := [...]string{"ComV1", node_id_type_component_tower,
		"HNODE", node_id_type_component_home,
		"ANODE", node_id_type_component_amplifier}
	match := false
	for _, code := range codes {
		if strings.Contains(id, code) {

			/* Check index of substring */
			idx := strings.Index(id, code)
			if idx == CODE_IDX {
				match = true
				break
			}

		}
	}

	if !match {
		err := errors.New("invalid Node Code")
		return "", err
	}

	/* RFC 1123 lowercase id and tags*/
	lid := strings.ToLower(id)

	return NodeID(lid), nil
}
