package ukama

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

const (
	CODE_IDX     = 10
	NodeIDLength = 23
	OEMCODE      = "UK"
	MFGCODE      = "SA"
	DELIMITER    = "-"
	HWVERSION    = "A1"
)

type NodeID string

func (n *NodeID) String() string{
	return string(*n)
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
	case "HomeNode", "homenode", "HOMENODE":
		code = "HNODE"
	case "CompNode", "compnode", "COMPNODE":
		code = "COMv1"
	case "AmpNode", "ampnode", "AMPNODE":
		code = "ANODE"
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
	weekstr := strconv.Itoa(week)
	code := GetNodeCodeForUnits(ntype)

	/*2+1+6+1+5+1+2+1+4*/
	/* UK-SA2154-HNODE-A1-XXXX*/
	uuid := OEMCODE + DELIMITER + MFGCODE + yearcode + weekstr + DELIMITER + code + DELIMITER + HWVERSION + DELIMITER + getRandCode(t)

	log.Infof("UUID: New NodeID for %s is %s and length is %d", ntype, uuid, len(uuid))

	/* RFC 1123 lowercase id and tags*/
	lid := strings.ToLower(uuid)

	return NodeID(lid)
}

func ValidateNodeId(id string) (NodeID, error) {

	/* TODO :: ADD more validation once we finalized this format */
	if len(id) != NodeIDLength {
		err := errors.New("invalid length")
		return "", err
	}

	/* Check for HW codes */
	codes := [...]string{"ComV1", "comv1", "HNODE", "hnode", "ANODE", "anode"}
	match := false
	for _, code := range codes {
		if strings.Contains(id, code) {

			/* Check index of sunstring */
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