package nmr

import (
	"fmt"

	"github.com/sirupsen/logrus"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/internal/order"
)

type NMR struct {
	S *sr.ServiceRouter
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func NewNMR(svcR *sr.ServiceRouter) *NMR {

	return &NMR{
		S: svcR,
	}
}

func (n *NMR) SendRestAPIReq(query map[string]string, body interface{}) error {

	errStatus := &ErrorMessage{}
	var err error
	var resp Response

	if body != nil {
		resp, err = n.S.C.R().
			SetError(errStatus).
			SetQueryParams(query).
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Get(n.S.Url.String() + "/service")
	} else {
		resp, err = n.S.C.R().
			SetError(errStatus).
			SetQueryParams(queryMap).
			Get(n.S.Url.String() + "/service")
	}

	if err != nil {
		logrus.Errorf("Failed to send api request to NMR. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("NMR update failure: %s", errStatus.Message)
	}

	return nil
}

func (n *NMR) NmrAddModule(module order.Module) {
	query := map[string]string{
		"module":module.ModuleID,
		"looking_to":"update_module",
	}

	err := n.SendRestAPIReq(query, node)
	if err != nil {
		logrus.Errorf("Failed to add module %s to NMR database.", module.ModuleID);
		return err
	}

	return nil
}

func (n *NMR) NmrDeleteModule() {

}

func (n *NMR) NmrUpdateModule() {

}

func (n *NMR) NmrUpdateStatus() {

}

func (n *NMR) NmrAssignModule() {

}

func (n *NMR) NmrAddNode(node order.Node) error {
	query := map[string]string{
		"node":node.NodeID,
		"looking_to":"update_node",
	}

	err := n.SendRestAPIReq(query, node)
	if err != nil {
		logrus.Errorf("Failed to add node %s to NMR database.", node.NodeID);
		return err
	}

	for idx,module := range node.Modules {
		err := n.AddModule(module)
		if err != nil {
			logrus.Errorf("Failed to add module %d  with id %s to NMR database.", idx, module.ModuleID);
			return err
		}

		logrus.Infof("Module %d with ID %s added to NMR database", idx. Node.Modules[idx].)
	}
}

func (n *NMR) NmrDeleteNode() {

}

func (n *NMR) NmrUpdateNode() {

}

func (n *NMR) NmrUpdateNodeStatus() {

}
