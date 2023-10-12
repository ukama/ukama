package configstore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"
	"google.golang.org/protobuf/types/known/anypb"

	utils "github.com/ukama/ukama/systems/node/configurator/pkg/utils"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

type ConfigStore struct {
	Store                providers.StoreProvider
	msgbus               mb.MsgBusServiceClient
	registrySystem       providers.RegistryProvider
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	configRepo           db.ConfigRepo
	commitRepo           db.CommitRepo
	OrgName              string
}

type ConfigStoreProvider interface {
	HandleConfigStoreEvent(ctx context.Context) error
	HandleConfigCommitReq(ctx context.Context, rVer string) error
}

type ConfigMetaData struct {
	org      string
	network  string
	site     string
	node     string
	app      string
	fileName string
}

const DIR_PREFIX = "/tmp/configstore/"
const PERM = 0755

func NewConfigStore(msgB mb.MsgBusServiceClient, registry providers.RegistryProvider, cfgDb db.ConfigRepo, cmtDb db.CommitRepo, orgName string, s providers.StoreProvider, t time.Duration) *ConfigStore {

	return &ConfigStore{
		Store:                s,
		registrySystem:       registry,
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

	// Get current commited version
	err = c.Store.GetRemoteConfigVersion(dir, cVerRec.Hash)
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	if lVer == cVerRec.Hash {
		log.Infof("HandleConfigStoreEvent remote config and current commit are same %s", cVerRec.Hash)
		return nil
	}

	return c.ProcessConfigStoreEvent(dir, cVerRec.Hash, lVer)
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

	// Get latest remote version
	err = c.Store.GetRemoteConfigVersion(dir, rVer)
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	cVerRec, err := c.commitRepo.GetLatest()
	if err != nil {
		log.Errorf("Failed to get latest commit: %v", err)
		return err
	}

	// Get current commited version
	err = c.Store.GetRemoteConfigVersion(dir, cVerRec.Hash)
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	if rVer == cVerRec.Hash {
		log.Infof("HandleConfigCommitReq remote config and requested commit are same %s", cVerRec.Hash)
		return nil
	}

	return c.ProcessConfigStoreEvent(dir, cVerRec.Hash, rVer)
}

func (c *ConfigStore) ProcessConfigStoreEvent(dir string, cVer string, rVer string) error {

	log.Infof("Looking for changes in config")

	filesUpdated, err := c.Store.GetDiff(cVer, rVer, dir+providers.LATEST_DIR_NAME)
	if err != nil {
		log.Errorf("Failed to get diff remote configs: %v", err)
		return err
	}

	var filesToUpdate []string
	cfPrefix := dir + providers.COMMIT_DIR_NAME + "/"
	lfPrefix := dir + providers.LATEST_DIR_NAME + "/"
	for _, file := range filesUpdated {
		_, change, err := utils.JsonDiff(cfPrefix+file, lfPrefix+file)
		if err != nil {
			log.Errorf("Failed to get json diff between %s and %s: %v", cfPrefix+file, lfPrefix+file, err)
			return err
		}

		if change {
			filesToUpdate = append(filesToUpdate, file)
		}
	}

	log.Infof("Files to be updated %+v", filesToUpdate)
	if len(filesToUpdate) > 0 {
		prepCommit := make(map[string]*pb.Config, len(filesToUpdate)) /* /* Map from file to config app, and real config files data*/
		prepNodeCommit := make(map[string][]string)                   /* Map from nodeId to config files*/
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
			cMetaData, err := ParseConfigStoreFilePath(file)
			if err != nil {
				log.Errorf("Failed to parse file %s. Error: %v", file, err)
				continue
			}

			/* This will filter out invalid network, site and nodes
			Also reads the file store the config */
			configToCommit, err := c.PrepareConfigCommit(cMetaData, lfPrefix+file)
			if err != nil {
				log.Errorf("Failed to prepare config commit for file %s and metadata %v. Error: %v", file, c, err)
				continue
			}

			prepMetaData[file] = cMetaData
			prepCommit[file] = configToCommit
			prepNodeCommit[cMetaData.node] = append(prepNodeCommit[cMetaData.node], file)
		}

		err = c.CommitConfig(prepCommit, prepNodeCommit, prepMetaData, rVer)
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

func (c *ConfigStore) PrepareConfigCommit(d *ConfigMetaData, file string) (*pb.Config, error) {

	log.Infof("Preparing commit config %s for node %+v", file, d)
	// var netId uuid.UUID
	// var err error

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

	data, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("unable to read file %s. Error %v", file, err)
		return nil, err
	}

	configReq := &pb.Config{
		Filename: filepath.Base(file), /* filename with path */
		App:      d.app,
		Data:     data,
	}

	return configReq, nil
}

func (c *ConfigStore) CommitConfig(m map[string]*pb.Config, nodes map[string][]string, md map[string]*ConfigMetaData, commit string) error {

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
		log.Infof("Pushing configs %+v for node %s", files, n)
		for _, f := range files {

			metaData = md[f]
			anyMsg, err := anypb.New(m[f])
			if err != nil {
				goto RecordState
			}

			msg := &pb.NodeFeederMessage{
				Target:     c.OrgName + "." + metaData.network + "." + metaData.site + "." + n,
				HTTPMethod: "POST",
				Path:       "/v1/configd/config",
				Msg:        anyMsg,
			}

			err = c.msgbus.PublishRequest(route, msg)
			if err != nil {
				log.Errorf("Failed to publish message %+v with key %+v. Errors %s", m[f], route, err.Error())
				goto RecordState
			}
			log.Infof("Published config %s on route %s for node %s ", msg, route, n)
			/* Atleast one is success */
			state = db.Partial
		}

		/* Publish config version information */
		err = c.PublishCommitInfo(metaData, route, commit)
		if err != nil {
			log.Errorf("Failed to pusblish the config version info.Erorr: %s", err.Error())
			goto RecordState
		}

		/* Publish the config commit info */
		state = db.Published

	RecordState:
		/* Update the version for commited config on node */
		cRec, err := c.configRepo.Get(n)
		if err != nil {
			log.Errorf("Failed to get last config for node %s.Error: %v", n, err)
			return err
		}

		cRec.Commit = db.Commit{Hash: commit}
		err = c.configRepo.UpdateLastCommit(*cRec, &state)
		if err != nil {
			log.Errorf("Failed to update latest commit: %v", err)
			return err
		}

	}

	return nil
}

func (c *ConfigStore) PublishCommitInfo(m *ConfigMetaData, route string, ver string) error {
	m.app = "config"
	m.fileName = "configInfo.json"

	type ConfigVerData struct {
		Name   string `json:"name"`
		App    string `json:"app"`
		Commit string `json:"commit"`
		NodeId string `json:"node_id"`
	}

	data := ConfigVerData{
		Name:   m.fileName,
		App:    m.app,
		Commit: ver,
		NodeId: m.node,
	}

	json, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Failed to marshal config version data %+v. Errors %s", data, err.Error())
		return err
	}

	configData := &pb.Config{
		Filename: m.fileName, /* filename with path */
		App:      m.app,
		Data:     json,
	}
	anyMsg, err := anypb.New(configData)
	if err != nil {
		return err
	}

	msg := &pb.NodeFeederMessage{
		Target:     c.OrgName + "." + m.network + "." + m.site + "." + m.node,
		HTTPMethod: "POST",
		Path:       "/v1/configd/config",
		Msg:        anyMsg,
	}

	err = c.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", configData, route, err.Error())
		return err
	}

	log.Infof("Published config %s on route %s for node %s ", msg, route, m.node)
	return nil
}
