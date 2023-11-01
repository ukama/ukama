/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/node/configurator/pkg/utils"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
)

type StoreProvider interface {
	GetLatestRemoteConfigs(dir string) (string, error)
	GetRemoteConfigVersion(dir string, version string) error
	GetDiff(prevSha string, curSha string, dir string) ([]string, error)
}

type gitClient struct {
	url  string
	user string
	pat  string
}

const LATEST_DIR_NAME = "/latest"
const COMMIT_DIR_NAME = "/commit"
const PERM = 0755

func (g *gitClient) GetLatestRemoteConfigs(dirPrefix string) (string, error) {

	err := utils.CreateDir(dirPrefix+LATEST_DIR_NAME, PERM)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	r, err := git.PlainClone(dirPrefix+LATEST_DIR_NAME, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: g.user, // yes, this can be anything except an empty string
			Password: g.pat,  //
		},
		URL: g.url,
	})

	if err != nil {
		return "", fmt.Errorf("error cloning config store from %s : %v", g.url, err)
	}

	ref, err := r.Head()
	if err != nil {
		return "", fmt.Errorf("error getting head revision %s : %v", g.url, err)
	}

	hash := ref.Hash()

	return hash.String(), nil
}

func (g *gitClient) GetRemoteConfigVersion(dirPrefix string, ver string) error {

	err := utils.CreateDir(dirPrefix+COMMIT_DIR_NAME, PERM)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	r, err := git.PlainClone(dirPrefix+COMMIT_DIR_NAME, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: g.user, // yes, this can be anything except an empty string
			Password: g.pat,  //
		},
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

func (g *gitClient) GetDiff(prevSha string, curSha string, dir string) ([]string, error) {

	var hash plumbing.Hash

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	hash = plumbing.NewHash(curSha)
	prevHash := plumbing.NewHash(prevSha)

	prevCommit, err := repo.CommitObject(prevHash)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	commit, err := repo.CommitObject(hash)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	log.Infof("Comparing from:" + prevCommit.Hash.String() + " to:" + commit.Hash.String())

	isAncestor, err := commit.IsAncestor(prevCommit)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	log.Infof("Is the prevCommit an ancestor of commit? : %v %v\n", commit, isAncestor)

	currentTree, err := commit.Tree()
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	prevTree, err := prevCommit.Tree()
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	patch, err := currentTree.Patch(prevTree)
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, err
	}

	log.Infof("Got %d changes", len(patch.Stats()))

	var changedFiles []string
	for _, fileStat := range patch.Stats() {
		log.Println(fileStat.Name)
		changedFiles = append(changedFiles, fileStat.Name)
	}
	return changedFiles, nil
}

func NewStoreClient(url string, user string, pat string, t time.Duration) (*gitClient, error) {

	N := &gitClient{
		url:  url,
		user: user,
		pat:  pat,
	}

	return N, nil
}
