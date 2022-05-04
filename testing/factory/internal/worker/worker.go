package worker

import (
	"fmt"

	"github.com/sirupsen/logrus"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/services/common/ukama"
	"github.com/ukama/ukama/testing/factory/internal"
	"github.com/ukama/ukama/testing/factory/internal/builder"
	"github.com/ukama/ukama/testing/factory/internal/nmr"
)

type Worker struct {
	b *builder.Build
	d *nmr.NMR
}

func NewWorker(r *sr.ServiceRouter) *Worker {
	fdb := nmr.NewNMR(r)

	return &Worker{
		d: fdb,
		b: builder.NewBuild(fdb),
	}

}

func (w *Worker) WorkerInit() {
	w.b.BuildInit()
}

/* Build Nodes */
func (w *Worker) WorkOnBuildOrder(ntype string, count int) ([]string, error) {

	idx := 0
	var nodeList []string

	for idx < count {

		var node internal.Node
		switch ntype {
		case ukama.NODE_ID_TYPE_HOMENODE, "hnode":
			node = NewHNode()
		case ukama.NODE_ID_TYPE_TOWERNODE, "tnode":
			node = NewTNode()
		case ukama.NODE_ID_TYPE_AMPNODE, "anode":
			node = NewANode()
		default:
			return nodeList, fmt.Errorf("unkown node type %s", ntype)
		}

		strId := string(node.NodeID)
		nodeList = append(nodeList, strId)

		/* Update the NMR DB */
		logrus.Debugf("Node %s is %+v", node.NodeID, node)
		err := w.d.NmrAddNode(node)
		if err != nil {
			/* TODO: May be collect errors for all node and then send response */
			logrus.Errorf("Failed to add node with nodeID %s. Error %s", node.NodeID, err.Error())
			return nodeList, fmt.Errorf("failed to add nodeID %s. Error %s", node.NodeID, err.Error())
		}

		/* Start bulding node */
		err = w.b.LaunchAndMonitorBuild(string(node.NodeID), node.Type)
		if err != nil {
			return nodeList, err
		}

		idx++

	}
	/* TODO check a way to rollback */
	return nodeList, nil
}
