package order

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/servcies/common/ukama"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/internal/builder"
	"github.com/ukama/ukama/testing/factory/internal/nmr"
)

type NodeBuilder struct {
	b  *builder.Build
	nc *nmr.NMR
}

func NewBuilder(r *sr.ServiceRouter) *NodeBuilder {
	return &NodeBuilder{
		b:  builder.NewBuild(),
		nc: nmr.NewNMR(r),
	}
}

/* Build Nodes */
func (nb *NodeBuilder) BuildNodes(ntype string, count int) error {

	idx := 0
	for idx < count {

		var node Node
		switch ntype {
		case ukama.ukama.NODE_ID_TYPE_HOMENODE:
			node = NewHNode()
		case ukama.NODE_ID_TYPE_TOWERNODE:
			node = NewHNode()
		case ukama.ukama.NODE_ID_TYPE_AMPLIFIERNODEconst:
			node = NewHNode()
		default:
			return fmt.Errorf("unkown node type %s", ntype)
		}

		/* Update the NMR DB */
		err := nb.nc.NmrAddNode(node)
		if err != nil {
			return err

		/* Start bulding node */
		err = nb.b.LaunchAndMonitorBuild(node.NodeID, node.Type)
		if err != nil {
			return err
		}

	}

	return nil
}