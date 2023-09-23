package providers

import (
	"fmt"
	"time"
	"ukama/ukama/systems/node/configurator/pkg/utils"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type StoreProvider interface {
	GetLatestRemoteConfigs() (string, string, error)
	GetRemoteConfigVersion(version string) (string, string, error)
}

type gitClient struct {
	url  string
	user string
	pat  string
}

const LATEST_DIR_NAME = "/tmp/configstore/latest"
const COMMIT_DIR_NAME = "/tmp/configstore/commit"
const PERM = 0755

func (g *gitClient) GetLatestRemoteConfigs() error {

	err := utils.CreateDir(LATEST_DIR_NAME, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	_, err = git.PlainClone(LATEST_DIR_NAME, false, &git.CloneOptions{
		URL: g.url,
	})

	if err != nil {
		return fmt.Errorf("error cloning config store from %s : %v", g.url, err)
	}

	return nil
}

func (g *gitClient) GetRemoteConfigVersion(ver string) error {

	err := utils.CreateDir(COMMIT_DIR_NAME, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	r, err := git.PlainClone(COMMIT_DIR_NAME, false, &git.CloneOptions{
		URL: g.url,
	})

	if err != nil {
		return fmt.Errorf("error cloning config store from %s : %v", g.url, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("error getting work tree: %v", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(ver),
	})
	if err != nil {
		return fmt.Errorf("error getting version %s: %v", ver, err)
	}
	return nil
}

func NewStoreClient(url string, user string, pat string, t time.Duration) (*gitClient, error) {

	N := &gitClient{
		url:  url,
		user: user,
		pat:  pat,
	}

	return N, nil
}
