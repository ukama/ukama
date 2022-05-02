package worker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/ukama"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/internal/builder"
	"github.com/ukama/ukama/testing/factory/internal/nmr"
	"github.com/ukama/ukama/testing/factory/internal"
)

type Worker struct {
	b  *builder.Build
	d *nmr.NMR
}

func NewWorker(r *sr.ServiceRouter) *Worker {
	
	return &Worker{
		d: nmr.NewNMR(r),
		b: builder.NewBuild(d),
	}
		
}

/* Build Nodes */
func (w *Worker) WorkOnBuildOrder(ntype string, count int) ([]string, error) {

	idx := 0
	nodeList := []string{}
	for idx < count {

		var node internal.Node
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

		append(nodeList, node.NodeID)
		
		/* Update the NMR DB */
		err := w.d.NmrAddNode(node)
		if err != nil {
			/* TODO: May be collect errors for all node and then send response */ 
			logrus.Errorf("Failed to add node with nodeID %s. Error %s", node.NodeId, err.Error())
			fmt.Errorf("failed to add nodeID %s. Error %s", err.Error())
			return nodeList, err
			
		/* Start bulding node */
		err = w.b.LaunchAndMonitorBuild(node.NodeID, node.Type)
		if err != nil {
			return nodeList, err
		}

		idx++

	}
	
	/* TODO check a way to rollback */
	return nodeList,nil
}