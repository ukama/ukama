/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package init

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"context"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

var config *pkg.Config

type InitData struct {
	OrgName                  string
	OrgIP                    string
	OrgCerts                 string
	SysName, SysIP, SysCerts string
	SysPort                  int32
	NodeId                   ukama.NodeID
	NodeIP, NodeCerts        string
	Init                     *InitClient
	Host                     string
	ROrgIp                   string
	MbHost                   string

	/* API requests */
	reqAddOrg    api.AddOrgRequest
	reqAddSystem api.AddSystemRequest
	reqAddNode   api.AddNodeRequest
	reqGetNode   api.GetNodeRequest
	reqGetSystem api.GetSystemRequest

	/* API responses */

}

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
}

func InitializeData() *InitData {
	config = pkg.NewConfig()

	d := &InitData{}
	d.OrgName = strings.ToLower(faker.FirstName()) + "_org"
	d.OrgIP = utils.RandomIPv4()
	d.OrgCerts = utils.RandomBase64String(2048)
	d.SysName = strings.ToLower(faker.FirstName()) + "_sys"
	d.SysIP = utils.RandomIPv4()
	d.SysCerts = utils.RandomBase64String(2048)
	d.SysPort = int32(utils.RandomPort())
	d.NodeIP = utils.RandomIPv4()
	d.NodeCerts = utils.RandomBase64String(2048)
	d.NodeId = ukama.NewVirtualHomeNodeId()
	d.Host = config.System.Init
	d.Init = NewInitClient(d.Host)
	d.MbHost = config.System.MessageBus

	d.reqAddOrg = api.AddOrgRequest{
		OrgName:     d.OrgName,
		Ip:          d.OrgIP,
		Certificate: d.OrgCerts,
	}

	d.reqAddSystem = api.AddSystemRequest{
		OrgName:     d.OrgName,
		SysName:     d.SysName,
		Ip:          d.SysIP,
		Certificate: d.SysCerts,
		Port:        int32(d.SysPort),
	}

	d.reqAddNode = api.AddNodeRequest{
		OrgName: d.OrgName,
		NodeId:  d.NodeId.String(),
	}

	d.reqGetNode = api.GetNodeRequest{
		OrgName: d.OrgName,
		NodeId:  d.NodeId.String(),
	}

	d.reqGetSystem = api.GetSystemRequest{
		OrgName: d.OrgName,
		SysName: d.SysName,
	}

	return d
}

var TC_init_add_org = &test.TestCase{
	Name:        "Add Organization.",
	Description: "Add organization to lookup table",
	Data:        &pb.AddOrgResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		a := tc.GetWorkflowData().(*InitData)
		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("Init").SetOrgName(a.OrgName).SetService("lookup").SetAction("create").SetObject("organization").MustBuild())
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Init.InitAddOrg(a.reqAddOrg)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		check := false
		d := tc.GetWorkflowData().(*InitData)
		resp := tc.GetData().(*pb.AddOrgResponse)
		if resp != nil && d.OrgName == resp.OrgName &&
			d.OrgIP == utils.IPv4CIDRToStringNotation(resp.Ip) &&
			d.OrgCerts == resp.Certificate &&
			tc.Watcher.Expections() {
			check = true
		}

		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		resp := tc.GetData().(*pb.AddOrgResponse)
		a := tc.GetWorkflowData().(*InitData)
		a.ROrgIp = resp.Ip
		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()
		return nil
	},
}

var TC_init_add_system = &test.TestCase{
	Name:        "Add System.",
	Description: "Add System to lookup table",
	Data:        &pb.AddSystemResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {

		a := tc.GetWorkflowData().(*InitData)
		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("Init").SetOrgName(a.OrgName).SetService("lookup").SetAction("create").SetObject("system").SetGlobalScope().MustBuild())
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Init.InitAddSystem(a.reqAddSystem)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		check := false
		d := tc.GetWorkflowData().(*InitData)
		resp := tc.GetData().(*pb.AddSystemResponse)
		if resp != nil && d.SysName == resp.SystemName &&
			d.SysIP == utils.IPv4CIDRToStringNotation(resp.Ip) &&
			d.SysCerts == resp.Certificate &&
			tc.Watcher.Expections() {
			check = true
		}
		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		tc.Watcher.Stop()
		return nil
	},
}

var TC_init_add_node = &test.TestCase{
	Name:        "Add Node.",
	Description: "Add node to a lookup table",
	Data:        &pb.AddSystemResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		a := tc.GetWorkflowData().(*InitData)
		log.Tracef("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem("Init").SetOrgName(a.OrgName).SetService("lookup").SetAction("create").SetObject("node").MustBuild())
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Init.InitAddNode(a.reqAddNode)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		check := false
		d := tc.GetWorkflowData().(*InitData)
		resp := tc.GetData().(*pb.AddNodeResponse)
		if resp != nil && d.NodeId.String() == resp.NodeId &&
			d.OrgName == resp.OrgName &&
			tc.Watcher.Expections() {
			check = true
		}
		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		tc.Watcher.Stop()
		return nil
	},
}

var TC_init_bootstrap_node = &test.TestCase{
	Name:        "Bootstrap Node.",
	Description: "Bootstrap node from a lookup table",
	Data:        &pb.AddSystemResponse{},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Init.InitGetNode(a.reqGetNode)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		check := false
		d := tc.GetWorkflowData().(*InitData)
		resp := tc.GetData().(*pb.GetNodeResponse)
		if resp != nil &&
			d.NodeId.String() == resp.NodeId &&
			d.OrgName == resp.OrgName {
			check = true
		}

		return check, nil
	},
}

var TC_init_get_system = &test.TestCase{

	Name:        "Get System.",
	Description: "get System from a lookup table",
	Data:        &pb.AddSystemResponse{},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Init.InitGetSystem(a.reqGetSystem)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		check := false
		d := tc.GetWorkflowData().(*InitData)
		resp := tc.GetData().(*pb.GetSystemResponse)
		if resp != nil && d.SysName == resp.SystemName &&
			d.SysIP == utils.IPv4CIDRToStringNotation(resp.Ip) &&
			d.SysCerts == resp.Certificate {
			check = true
		}

		return check, nil
	},
}
