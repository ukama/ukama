/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package configstore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	utils "github.com/ukama/ukama/systems/node/configurator/pkg/utils"
)

type ConfigStore struct {
	Store                providers.StoreProvider
	msgbus               mb.MsgBusServiceClient
	networkClient        creg.NetworkClient
	siteClient           creg.SiteClient
	nodeClient           creg.NodeClient
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	configRepo           db.ConfigRepo
	commitRepo           db.CommitRepo
	OrgName              string
}

type FilesToUpdate struct {
	Name   string
	Reason int
}

type ConfigStoreProvider interface {
	HandleConfigStoreEvent(ctx context.Context) error
	HandleConfigCommitReq(ctx context.Context, rVer string) error
	HandleConfigCommitReqForNode(ctx context.Context, rVer string, nodeid string) error
}

const (
	REASON_UNKNOWN = iota
	REASON_ADDED   = 1
	REASON_DELETED = 2
	REASON_UPDATED = 3
)

type ConfigData struct {
	FileName  string `json:"file_name"`
	App       string `json:"app"`
	Version   string `json:"version"`
	Data      []byte `json:"data"`
	Reason    int    `json:"reason"`
	Timestamp uint32 `json:"timestamp"`
	FileCount int    `json:"file_count"`
}

type ConfigMetaData struct {
	network  string
	site     string
	node     string
	app      string
	fileName string
}

const DIR_PREFIX = "/tmp/configstore/"
const PERM = 0755

func NewConfigStore(msgB mb.MsgBusServiceClient, cnet creg.NetworkClient, csite creg.SiteClient, cnode creg.NodeClient, cfgDb db.ConfigRepo, cmtDb db.CommitRepo, orgName string, s providers.StoreProvider, t time.Duration) *ConfigStore {

	return &ConfigStore{
		Store:                s,
		networkClient:        cnet,
		siteClient:           csite,
		nodeClient:           cnode,
		msgbus:               msgB,
		NodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetRequestType().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName), //Need to have something same to other routes
		OrgName:              orgName,
		configRepo:           cfgDb,
		commitRepo:           cmtDb,
	}
}

func (c *ConfigStore) HandleConfigStoreEvent(ctx context.Context) error {
	log.Infof("HandleConfigStoreEvent")

	dir := DIR_PREFIX + utils.RandomDirName()

	err := utils.CreateDir(dir, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	defer func() {
		err := utils.RemoveDir(dir)
		if err != nil {
			log.Errorf("error removing directory: %v", err)
		}
	}()

	// Get latest remote version
	lVer, err := c.Store.GetLatestRemoteConfigs(dir)
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	/* Get current commit */
	cVerRec, err := c.commitRepo.GetLatest()
	if err != nil {
		log.Errorf("Failed to get latest commit: %v", err)
		return err
	}

	// Get current committed version
	err = c.Store.GetRemoteConfigVersion(dir, cVerRec.Hash)
	if err != nil {
		log.Errorf("Failed to get current remote configs: %v", err)
		return err
	}

	if lVer == cVerRec.Hash {
		log.Infof("HandleConfigStoreEvent remote config and current commit are same %s", cVerRec.Hash)
		return nil
	}

	files, dir, err := c.LookingForChanges(dir, cVerRec.Hash, lVer)
	if err != nil {
		log.Errorf("Failed to get change list for version %s from version %s.", lVer, cVerRec.Hash)
		return err
	}

	return c.ProcessConfigStoreEvent(files, lVer, dir)

}

func (c *ConfigStore) HandleConfigCommitReq(ctx context.Context, rVer string) error {
	log.Infof("HandleConfigCommitReq")

	dir := DIR_PREFIX + utils.RandomDirName()

	err := utils.CreateDir(dir, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	defer func() {
		err := utils.RemoveDir(dir)
		if err != nil {
			log.Errorf("error removing directory: %v", err)
		}
	}()

	// Get requested remote version
	err = c.Store.GetRemoteConfigVersion(dir, rVer)
	if err != nil {
		log.Errorf("Failed to get requested remote configs: %v", err)
		return err
	}

	cVerRec, err := c.commitRepo.GetLatest()
	if err != nil {
		log.Errorf("Failed to get latest commit: %v", err)
		return err
	}

	if rVer == cVerRec.Hash {
		log.Infof("HandleConfigCommitReq remote config and requested commit are same %s", cVerRec.Hash)
		return nil
	}

	// Get current committed version
	err = c.Store.GetRemoteConfigVersion(dir, cVerRec.Hash)
	if err != nil {
		log.Errorf("Failed to get current remote configs: %v", err)
		return err
	}

	files, dir, err := c.LookingForChanges(dir, cVerRec.Hash, rVer)
	if err != nil {
		log.Errorf("Failed to get change list for version %s from version %s.", rVer, cVerRec.Hash)
		return err
	}

	return c.ProcessConfigStoreEvent(files, rVer, dir)
}

func (c *ConfigStore) HandleConfigCommitReqForNode(ctx context.Context, rVer string, nodeid string) error {
	log.Infof("HandleConfigCommitReqForNode")

	dir := DIR_PREFIX + utils.RandomDirName()

	err := utils.CreateDir(dir, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	defer func() {
		err := utils.RemoveDir(dir)
		if err != nil {
			log.Errorf("error removing directory: %v", err)
		}
	}()

	// Get requested remote version
	err = c.Store.GetRemoteConfigVersion(dir, rVer)
	if err != nil {
		log.Errorf("Failed to get requested remote configs: %v", err)
		return err
	}

	files, dir, err := c.LookingForNodeConfigs(dir, nodeid, rVer)
	if err != nil {
		log.Errorf("Failed to get change list of configs for node %s in rev %s.", nodeid, rVer)
		return err
	}

	return c.ProcessConfigStoreEvent(files, rVer, dir)
}

func (c *ConfigStore) LookingForNodeConfigs(dir string, nodeId string, rVer string) ([]FilesToUpdate, string, error) {
	log.Infof("Looking for nodeid %s configs", nodeId)

	path, err := utils.FindDir(nodeId, dir+providers.COMMIT_DIR_NAME)
	if err != nil {
		log.Errorf("Failed to find nodeid %s config under %s dir", nodeId, dir+providers.COMMIT_DIR_NAME)
		return nil, "", err
	}

	filesUpdated, err := utils.GetFiles(*path)
	if err != nil {
		log.Errorf("Failed to get diff remote configs: %v", err)
		return nil, "", err
	}

	var filesToUpdate []FilesToUpdate
	lprefix := dir + providers.COMMIT_DIR_NAME + "/"
	for _, file := range filesUpdated {
		filePath := strings.Split(file, lprefix)
		filesToUpdate = append(filesToUpdate, FilesToUpdate{Name: filePath[1], Reason: REASON_ADDED})
	}

	log.Infof("Files to be updated %+v", filesToUpdate)
	return filesToUpdate, lprefix, nil
}

func (c *ConfigStore) LookingForChanges(dir string, cVer string, rVer string) ([]FilesToUpdate, string, error) {
	log.Infof("Looking for changes in config")

	filesUpdated, err := c.Store.GetDiff(cVer, rVer, dir+providers.LATEST_DIR_NAME)
	if err != nil {
		log.Errorf("Failed to get diff remote configs: %v", err)
		return nil, "", err
	}

	var filesToUpdate []FilesToUpdate
	cfPrefix := dir + providers.COMMIT_DIR_NAME + "/"
	lfPrefix := dir + providers.LATEST_DIR_NAME + "/"
	for _, file := range filesUpdated {
		_, change, reason, err := utils.JsonDiff(cfPrefix+file, lfPrefix+file)
		if err != nil {
			log.Errorf("Failed to get json diff between %s and %s: %v", cfPrefix+file, lfPrefix+file, err)
			return nil, "", err
		}

		if change {
			filesToUpdate = append(filesToUpdate, FilesToUpdate{Name: file, Reason: reason})
		}
	}

	log.Infof("Files to be updated %+v", filesToUpdate)
	return filesToUpdate, lfPrefix, nil
}

func (c *ConfigStore) ProcessConfigStoreEvent(filesToUpdate []FilesToUpdate, rVer string, dir string) error {

	if len(filesToUpdate) > 0 {
		prepCommit := make(map[string]*ConfigData, len(filesToUpdate)) /* /* Map from file to config app, and real config files data*/
		prepNodeCommit := make(map[string][]string)                    /* Map from nodeId to config files*/
		prepMetaData := make(map[string]*ConfigMetaData)
		for _, file := range filesToUpdate {
			/* Get the meta information about config from the path of the filename
			  networkABC/site123/uk-sa1000-HNODE-2145/epc/sctp.json
				Org:ukama
				Network: networkABC
				Site:site123
				Node: uk-sa1000-HNODE-2145
				App: epc
			*/
			cMetaData, err := ParseConfigStoreFilePath(file.Name)
			if err != nil {
				log.Errorf("Failed to parse file %s. Error: %v", file.Name, err)
				continue
			}

			/* This will filter out invalid network, site and nodes
			Also reads the file store the config */
			configToCommit, err := c.PrepareConfigCommit(cMetaData, dir+file.Name, file.Reason)
			if err != nil {
				log.Errorf("Failed to prepare config commit for file %s and metadata %v. Error: %v", file.Name, c, err)
				continue
			}
			configToCommit.Reason = file.Reason
			configToCommit.Version = rVer

			prepMetaData[file.Name] = cMetaData
			prepCommit[file.Name] = configToCommit
			prepNodeCommit[cMetaData.node] = append(prepNodeCommit[cMetaData.node], file.Name)
		}

		err := c.CommitConfig(prepCommit, prepNodeCommit, prepMetaData, rVer)
		if err != nil {
			return err
		}

	} else {
		log.Info("No changes to commit.")
	}

	return nil
}

/* Parse ukama/networkABC/site123/uk-sa1000-HNODE-2145/epc/sctp.json */
func ParseConfigStoreFilePath(path string) (*ConfigMetaData, error) {

	c := &ConfigMetaData{}
	p := strings.Split(path, "/")
	log.Infof("Creating metadata for file path %s {%+v}", path, p)
	fnPos := len(p) - 1

	if fnPos > 5 {
		log.Errorf("Invalid path length %s", path)
		return nil, fmt.Errorf("invalid path length %s", path)
	}

	fn := p[fnPos]
	if !IfFileName(fn) {
		log.Errorf("Invalid path for config %s", path)
		return nil, fmt.Errorf("invalid path for config file %s", path)
	}

	for i, pe := range p {
		switch i {
		case 0:
			c.network = pe
		case 1:
			c.site = pe
		case 2:
			c.node = pe
		case 3:
			c.app = pe
		case 4:
			c.fileName = pe
		default:
			return nil, fmt.Errorf("invalid path element at %d of %s", i, path)
		}
	}

	return c, nil

}

var ExpectedConfigExt = []string{"json"}

func IfFileName(f string) bool {
	fileName := false
	fp := strings.Split(f, ".")
	fe := fp[len(fp)-1]
	for _, e := range ExpectedConfigExt {
		if fe == e {
			fileName = true
			break
		}
	}
	return fileName
}

func (c *ConfigStore) PrepareConfigCommit(d *ConfigMetaData, file string, reason int) (*ConfigData, error) {

	log.Infof("Preparing commit config %s for node %+v", file, d)
	// var netId uuid.UUID
	var err error

	// netId, err = uuid.FromString(d.network)
	// if err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid network ID format: %s", err.Error())
	// }

	// if err := c.registrySystem.ValidateNetwork(netId.String(), d.org); err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid network ID: %s", err.Error())
	// }

	// if err := c.registrySystem.ValidateSite(d.network, d.site, d.org); err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid site %s", err.Error())
	// }

	// if err := c.registrySystem.ValidateNode(d.node, d.org); err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "invalid node: %s", err.Error())
	// }

	var data []byte
	if reason != REASON_DELETED {
		data, err = os.ReadFile(file)
		if err != nil {
			log.Errorf("unable to read file %s. Error %v", file, err)
			return nil, err
		}
	}

	configReq := &ConfigData{
		FileName: filepath.Base(file), /* filename with path */
		App:      d.app,
		Data:     data,
	}

	return configReq, nil
}

func (c *ConfigStore) CommitConfig(m map[string]*ConfigData, nodes map[string][]string, md map[string]*ConfigMetaData, commit string) error {

	route := c.NodeFeederRoutingKey.SetObject("node").SetAction("publish").MustBuild()

	for n, files := range nodes {
		state := db.Failed

		/* Check if node existes in configuration db */
		_, err := c.configRepo.Get(n)
		if err != nil {
			log.Errorf("Failed to get configuration data for node %s: %v", n, err)
			continue
		}

		metaData := &ConfigMetaData{}
		t := (uint32)(time.Now().Unix())
		count := len(files) + 1
		log.Infof("Pushing configs %+v for node %s with timestamp %d", files, n, t)
		for _, f := range files {

			metaData = md[f]

			cd := m[f]
			cd.Timestamp = t
			cd.FileCount = count

			jd, err := json.Marshal(cd)
			if err != nil {
				log.Errorf("Failed to marshal configdata %+v. Errors %s", cd, err.Error())
				return err
			}

			msg := &pb.NodeFeederMessage{
				Target:     c.OrgName + "." + metaData.network + "." + metaData.site + "." + n,
				HTTPMethod: "POST",
				Path:       "configd/v1/config",
				Msg:        jd,
			}

			err = c.msgbus.PublishRequest(route, msg)
			if err != nil {
				log.Errorf("Failed to publish message %+v with key %+v. Errors %s", m[f], route, err.Error())
				goto RecordState
			}
			log.Infof("Published config %s  with timestamp %d on route %s for node %s ", msg, t, route, n)

			/* Atleast one is success */
			state = db.Partial
		}

		/* Publish config version information */
		err = c.PublishCommitInfo(metaData, route, commit, t, count)
		if err != nil {
			log.Errorf("Failed to pusblish the config version info.Erorr: %s", err.Error())
			goto RecordState
		}

		/* Publish the config commit info */
		state = db.Published

	RecordState:
		/* Update the version for committed config on node */
		cRec, err := c.configRepo.Get(n)
		if err != nil {
			log.Errorf("Failed to get last config for node %s.Error: %v", n, err)
			return err
		}

		cRec.Commit = db.Commit{Hash: commit}

		/* Get commit */
		cmt, err := c.commitRepo.Get(commit)
		if err == nil {
			cRec.Commit = *cmt
		}

		err = c.configRepo.UpdateLastCommit(*cRec, &state)
		if err != nil {
			log.Errorf("Failed to update latest commit: %v", err)
			return err
		}

	}

	return nil
}

func (c *ConfigStore) PublishCommitInfo(m *ConfigMetaData, route string, ver string, t uint32, count int) error {
	m.app = "configd"
	m.fileName = "version.json"

	cdata := ConfigData{
		FileName:  m.fileName, /* filename with path */
		App:       m.app,
		Data:      []byte(""),
		Reason:    REASON_UPDATED,
		Timestamp: t,
		Version:   ver,
		FileCount: count,
	}

	jd, err := json.Marshal(&cdata)
	if err != nil {
		log.Errorf("Failed to marshal config version data %+v. Errors %s", cdata, err.Error())
		return err
	}

	cdata.Data = jd
	jsonMsg, err := json.Marshal(&cdata)
	if err != nil {
		log.Errorf("Failed to marshal configdata %+v. Errors %s", cdata, err.Error())
		return err
	}

	msg := &pb.NodeFeederMessage{
		Target:     c.OrgName + "." + m.network + "." + m.site + "." + m.node,
		HTTPMethod: "POST",
		Path:       "configd/v1/config",
		Msg:        jsonMsg,
	}

	err = c.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", jsonMsg, route, err.Error())
		return err
	}

	log.Infof("Published config %s on route %s for node %s ", msg, route, m.node)
	return nil
}
