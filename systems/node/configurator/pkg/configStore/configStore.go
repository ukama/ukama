package config

import (
	"context"
	"strings"
	"time"
	"ukama/ukama/systems/node/configurator/pkg/providers"

	utils "command-line-arguments/home/vishal/cdrive/work/git/ukama/ukama/systems/node/configurator/pkg/utils/json.go"

	log "github.com/sirupsen/logrus"
)

type ConfigStore struct {
	Store providers.StoreProvider
}

func NewConfigStore(url string, user string, pat string, t time.Duration) *ConfigStore {
	s, err := providers.NewStoreClient(url, user, pat, t)
	if err != nil {
		return nil
	}

	return &ConfigStore{
		Store: s,
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
		_, change, err := utils.JsonDiff(strings.Join(providers.COMMIT_DIR_NAME, srcFile), strings.Join(providers.COMMIT_DIR_NAME, targetFile))
		if change {
			filesToUpdate = append(filesToUpdate, file)
		}
	}

	/* Get the meta information about config from the path of the filename
	  ukama/networkABC/site123/uk-sa1000-HNODE-2145/epc/sctp.json
		Org:ukama
		Network: networkABC
		Site:site123
		Node: uk-sa1000-HNODE-2145
		App: epc
	*/

	/* Send these files  to the device */

	/* Update the current version for the node */

}
