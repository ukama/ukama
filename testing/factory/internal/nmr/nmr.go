package nmr

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/internal"
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

func (n *NMR) SendRestAPIReq(query map[string]string, body ...interface{}) error {

	errStatus := &ErrorMessage{}
	var err error
	resp := &resty.Response{}
	for _, item := range body {
		resp, err = n.S.C.R().
			SetError(errStatus).
			SetQueryParams(query).
			SetHeader("Content-Type", "application/json").
			SetBody(item).
			Get(n.S.Url.String() + "/service")
	}

	if len(body) == 0 {
		resp, err = n.S.C.R().
			SetError(errStatus).
			SetQueryParams(query).
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

func (n *NMR) NmrAddModule(module internal.Module) error {
	query := map[string]string{
		"module":     string(module.ModuleID),
		"looking_to": "update_module",
	}

	err := n.SendRestAPIReq(query, module)
	if err != nil {
		logrus.Errorf("Failed to add module %s to NMR database.", module.ModuleID)
		return err
	}

	return nil
}

func (n *NMR) NmrAssignModule(nodeID string, moduleID string) error {
	query := map[string]string{
		"node":       nodeID,
		"looking_to": "allocate",
		"module":     moduleID,
	}

	err := n.SendRestAPIReq(query, nil)
	if err != nil {
		logrus.Errorf("Failed to allocate module %s to  nodeID %s status in NMR database. Error: %s", moduleID, nodeID, err.Error())
		return err
	}

	logrus.Info("Status allocated moduleID %s for NodeID: %s in NMR Database.", moduleID, nodeID)
	return nil

}

func (n *NMR) NmrAddNode(node internal.Node) error {
	query := map[string]string{
		"node":       string(node.NodeID),
		"looking_to": "update_node",
	}

	err := n.SendRestAPIReq(query, node)
	if err != nil {
		logrus.Errorf("Failed to add node %s to NMR database.", node.NodeID)
		//return err
	}

	for idx, module := range node.Modules {
		err := n.NmrAddModule(module)
		if err != nil {
			logrus.Errorf("Failed to add module %d  with id %s to NMR database.", idx, module.ModuleID)
			//return err
		}
		logrus.Infof("Module %d with ID %s added to NMR database", idx, node.Modules[idx].ModuleID)

		/* Allocate Module */
		err = n.NmrAssignModule(string(node.NodeID), string(node.Modules[idx].ModuleID))
		if err != nil {
			logrus.Errorf("Failed to add module %d  with id %s to NMR database.", idx, module.ModuleID)
			//return err
		}
		logrus.Infof("Module %d with ID %s added to NMR database", idx, node.Modules[idx].ModuleID)
	}

	return nil
}

func (n *NMR) NmrUpdateNodeStatus(nodeID string, status string) error {
	query := map[string]string{
		"node":       nodeID,
		"looking_to": "status_update",
		"status":     status,
	}

	err := n.SendRestAPIReq(query, nil)
	if err != nil {
		logrus.Errorf("Failed to update nodeID %s status in NMR database. Error: %s", nodeID, err.Error())
		return err
	}

	logrus.Info("Status updated for NodeID: %s with %s in NMR Database.", nodeID, status)
	return nil
}
