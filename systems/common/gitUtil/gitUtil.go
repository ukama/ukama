/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package gitUtil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	thttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Company struct {
	Company       string `json:"company"`
	GitBranchName string `json:"git_branch_name"`
	Email         string `json:"email"`
}

type Environment struct {
	Production []Company `json:"production"`
	Test       []Company `json:"test"`
}

type Component struct {
	Company       string `json:"company"`
	Category      string `json:"category"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer"`
	Managed       string `json:"managed"`
	Warranty      uint32 `json:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}

func httpRequest(token string, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func RemoveTempDir(path string) {
	os.RemoveAll("temp")
}

func CreateTempDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func RemoveTempDirIfExist(path string) {
	RemoveTempDir(path)
	CreateTempDir(path)
}

func ReadRootFile(token string, url string) (*Environment, error) {
	body, err := httpRequest(token, url)
	if err != nil {
		return nil, err
	}

	var env Environment
	err = json.Unmarshal(body, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func CloneGitRepo(url string, username string, token string, path string) (*git.Repository, error) {
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &thttp.BasicAuth{
			Username: username,
			Password: token,
		},

		SingleBranch: true,
		URL:          url,
		Progress:     os.Stdout,
	})

	return r, err
}

func BranchCheckout(r *git.Repository, branch string, username string, token string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}

	if err := w.Checkout(&branchCoOpts); err != nil {
		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		err = fetchOrigin(r, mirrorRemoteBranchRefSpec, username, token)
		if err != nil {
			return err
		}

		err = w.Checkout(&branchCoOpts)
		if err != nil {
			return err
		}

		log.Printf("branch checkout success %s", branch)
	}

	return nil
}

func fetchOrigin(repo *git.Repository, refSpecStr string, username string, token string) error {
	remotes, err := repo.Remotes()
	if err != nil {
		return err
	}

	remote, err := repo.Remote(remotes[0].Config().Name)
	if err != nil {
		return err
	}

	if err = remote.Fetch(&git.FetchOptions{
		Auth: &thttp.BasicAuth{
			Username: username,
			Password: token,
		},
		RemoteURL: remote.Config().URLs[0],
		RefSpecs: []config.RefSpec{
			config.RefSpec(refSpecStr),
		},
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

func GetFilesPath(lpath string) (map[string][]string, error) {
	set := make(map[string][]string)
	set["accounting"] = []string{}
	set["components"] = []string{}
	set["contracts"] = []string{}
	err := filepath.Walk(lpath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				rp := strings.Replace(path, lpath, "", -1)
				parts := strings.Split(rp, "/")
				if len(parts) > 2 && (parts[0] == "components" || parts[0] == "contracts") {
					set[parts[0]] = append(set[parts[0]], path)
				} else if len(parts) > 1 && parts[0] == "accounting" {
					set[parts[0]] = append(set[parts[0]], path)
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return set, nil
}

func ReadFile(path string) (map[string]interface{}, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var jsonObj map[string]interface{}
	err = json.Unmarshal(content, &jsonObj)
	if err != nil {
		return nil, err
	}

	return jsonObj, nil
}
