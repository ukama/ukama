package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/orchestrator/constructor/pkg"
	"github.com/ukama/ukama/systems/common/msgbus"

	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"

	utils "github.com/ukama/ukama/systems/node/configurator/pkg/utils"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfigStore struct {
	Store                providers.StoreProvider
	msgbus               mb.MsgBusServiceClient
	registrySystem       providers.RegistryProvider
	NodeFeederRoutingKey msgbus.RoutingKeyBuilder
	OrgName              string
}

type ConfigMetaData struct {
	org      string
	network  string
	site     string
	node     string
	app      string
	fileName string
}

func NewConfigStore(msgB mb.MsgBusServiceClient, registry providers.RegistryProvider, orgName string, url string, user string, pat string, t time.Duration) *ConfigStore {
	s, err := providers.NewStoreClient(url, user, pat, t)
	if err != nil {
		return nil
	}

	return &ConfigStore{
		Store:                s,
		registrySystem:       registry,
		msgbus:               msgB,
		NodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		OrgName:              orgName,
	}
}

func (c *ConfigStore) HandleConfigStoreEvent(ctx context.Context) error {
	log.Infof("HandleConfigStoreEvent")

	// Get latest remote version
	lv, err := c.Store.GetLatestRemoteConfigs()
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	/* TODO: Get current commit */
	currentCommit := ""

	// Get current commited version
	err = c.Store.GetRemoteConfigVersion(currentCommit)
	if err != nil {
		log.Errorf("Failed to get latest remote configs: %v", err)
		return err
	}

	filesUpdated, err := c.Store.GetDiff(currentCommit, lv, providers.LATEST_DIR_NAME)
	if err != nil {
		log.Errorf("Failed to get diff remote configs: %v", err)
		return err
	}

	var filesToUpdate []string
	for _, file := range filesUpdated {
		_, change, err := utils.JsonDiff(providers.COMMIT_DIR_NAME+file, providers.COMMIT_DIR_NAME+file)
		if err != nil {
			log.Errorf("Failed to get diff between %s and %s: %v", providers.COMMIT_DIR_NAME+file, providers.COMMIT_DIR_NAME+file, err)
		}
		if change {
			filesToUpdate = append(filesToUpdate, file)
		}
	}

	prepCommit := make(map[string]*pb.Config, len(filesToUpdate))
	for _, file := range filesUpdated {
		/* Get the meta information about config from the path of the filename
		  ukama/networkABC/site123/uk-sa1000-HNODE-2145/epc/sctp.json
			Org:ukama
			Network: networkABC
			Site:site123
			Node: uk-sa1000-HNODE-2145
			App: epc
		*/
		cMetaData, err := ParseConfigStoreFilePath(file)
		if err != nil {
			return err
		}

		configToCommit, err := c.PrepareConfigCommit(cMetaData, file)
		if err != nil {
			log.Errorf("Failed to prepare config commit for file %s and metadata %v. Error: %v", file, c, err)
			return err
		}
		prepCommit[file] = configToCommit
	}

	err = c.CommitConfig(prepCommit)
	if err != nil {
		return err
	}

	/* Update the version for commited config */

	return nil
}

/* Parse ukama/networkABC/site123/uk-sa1000-HNODE-2145/epc/sctp.json */
func ParseConfigStoreFilePath(path string) (*ConfigMetaData, error) {
	c := &ConfigMetaData{}
	p := strings.Split(path, "/")
	fnPos := len(p) - 1
	fn := p[fnPos]
	if !IfFileName(fn) {
		log.Errorf("Invalid path for config %s", path)
		return nil, fmt.Errorf("invalid path for config file %s", path)
	}

	for i, pe := range p {
		switch i {
		case 0:
			c.org = pe
		case 1:
			c.network = pe
		case 2:
			c.site = pe
		case 3:
			c.node = pe
		case 4:
			c.app = pe
		case 5:
			c.app = pe
		default:
			return nil, fmt.Errorf("invalid path element %s", path)
		}
	}

	return c, nil

}

var ExpectedConfigExt = []string{"json"}

func IfFileName(f string) bool {
	fileName := false
	fp := strings.Split(f, ".")
	fe := fp[len(f)-1]
	for _, e := range ExpectedConfigExt {
		if fe == e {
			fileName = true
			break
		}
	}
	return fileName
}

func (c *ConfigStore) PrepareConfigCommit(d *ConfigMetaData, file string) (*pb.Config, error) {

	log.Infof("Sending config %s to node %+v", file, c)
	var netId uuid.UUID
	var err error
	if d.network == "" {
		d.network = "*"
	} else {

		netId, err = uuid.FromString(d.network)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid network ID format: %s", err.Error())
		}

		if err := c.registrySystem.ValidateNetwork(netId.String(), d.org); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid network ID: %s", err.Error())
		}
	}

	if d.site == "" {
		d.site = "*"
	} else {
		if err := c.registrySystem.ValidateSite(d.network, d.site, d.org); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid site %s", err.Error())
		}
	}

	if d.node == "" {
		d.node = "*"
	} else {
		if err := c.registrySystem.ValidateNode(d.node, d.org); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid node: %s", err.Error())
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("unable to read file %s. Error %v", file, err)
		return nil, err
	}

	configReq := &pb.Config{
		Filename: file, /* filename with path */
		App:      d.app,
		Data:     data,
	}

	return configReq, nil
}

func (c *ConfigStore) CommitConfig(m map[string]*pb.Config) error {
	route := c.NodeFeederRoutingKey.SetActionUpdate().SetObject("config").MustBuild()

	for _, d := range m {
		err := c.msgbus.PublishRequest(route, d)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", d, route, err.Error())
			return err
		}
		log.Infof("Published commit for route %s and config %s", route, d.Filename)

	}

	return nil
}
