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
	NODE_ID_TYPE_HOMENODE  = "hnode"
	NODE_ID_TYPE_TOWERNODE = "tnode"
	NODE_ID_TYPE_AMPNODE   = "anode"
	NODE_ID_TYPE_UNDEFINED = "undef"
)

const (
	MODULE_ID_TYPE_COMP      = "comv1"
	MODULE_ID_TYPE_TRX       = "trx"
	MODULE_ID_TYPE_CTRL      = "ctrl"
	MODULE_ID_TYPE_FE        = "fe"
	MODULE_ID_TYPE_UNDEFINED = "undef"
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

	case node_id_type_component_amplifier:
		return NODE_ID_TYPE_AMPNODE

	case node_id_type_component_tower:
		return NODE_ID_TYPE_TOWERNODE
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
	switch strings.ToLower(ntype) {
	case NODE_ID_TYPE_HOMENODE, "homenode":
		code = NODE_ID_TYPE_HOMENODE
	case NODE_ID_TYPE_TOWERNODE, "towernode":
		code = NODE_ID_TYPE_TOWERNODE
	case NODE_ID_TYPE_AMPNODE, "ampnode":
		code = NODE_ID_TYPE_AMPNODE
	default:
		code = NODE_ID_TYPE_UNDEFINED
	}
	return code
}

/* Get HW Code */
func GetModuleCodeForUnits(mtype string) string {
	var code string
	switch strings.ToLower(mtype) {
	case MODULE_ID_TYPE_COMP, "com":
		code = MODULE_ID_TYPE_COMP
	case MODULE_ID_TYPE_TRX, "transciever":
		code = MODULE_ID_TYPE_TRX
	case MODULE_ID_TYPE_CTRL, "control":
		code = MODULE_ID_TYPE_CTRL
	case MODULE_ID_TYPE_FE, "rffe", "frontend":
		code = MODULE_ID_TYPE_FE
	default:
		code = MODULE_ID_TYPE_UNDEFINED
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
	codes := [...]string{
		node_id_type_component_home,
		node_id_type_component_amplifier,
		node_id_type_component_tower}
	match := false
	for _, code := range codes {
		if strings.Contains(strings.ToLower(id), code) {

			/* Check index of substring */
			idx := strings.Index(strings.ToLower(id), code)
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
