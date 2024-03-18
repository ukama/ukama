/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	git "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/services/gitClient/pb/gen"
)

const path = "temp/repo/networks/"

type GitClientServer struct {
	url           string
	token         string
	username      string
	rootConfigUrl string
	r             *git.Repository
	pb.UnimplementedGitClientServiceServer
}

func NewGitClientServer(url string, username string, token string, rootConfigUrl string) *GitClientServer {
	return &GitClientServer{
		url:           url,
		token:         token,
		username:      username,
		rootConfigUrl: rootConfigUrl,
	}
}

func (g *GitClientServer) FetchComponents(ctx context.Context, req *pb.FetchComponentsRequest) (*pb.FetchComponentsResponse, error) {
	log.Infof("Fetching componente")

	RemoveTempDirIfExist(path)
	r, err := CloneGitRepo(g.url, g.username, g.token, path)
	CheckIfError(err)
	g.r = r

	components := []Component{}
	res, _ := ReadRootFile(g.token, g.rootConfigUrl)
	for _, company := range res.Test {
		BranchCheckout(r, company.GitBranchName, g.username, g.token)
		paths, _ := GetFilesPath(path)
		for _, p := range paths["components"] {
			component, err := ReadFile(p, company.Company)
			components = append(components, component)
			CheckIfError(err)
		}
	}

	RemoveTempDir(path)

	return &pb.FetchComponentsResponse{Component: ComponentsToPbComponents(components)}, nil
}

func JSONToPb(components Component) *pb.Component {
	return &pb.Component{
		Company:       components.Company,
		InventoryId:   components.InventoryID,
		Category:      components.Category,
		Type:          pb.ComponentType(pb.ComponentType_value[components.Type]),
		Description:   components.Description,
		DatasheetURL:  components.DatasheetURL,
		ImagesURL:     components.ImagesURL,
		PartNumber:    components.PartNumber,
		Manufacturer:  components.Manufacturer,
		Managed:       components.Managed,
		Warranty:      components.Warranty,
		Specification: components.Specification,
	}
}

func ComponentsToPbComponents(components []Component) []*pb.Component {
	var pbComponents []*pb.Component
	for _, c := range components {
		pbComponents = append(pbComponents, JSONToPb(c))
	}
	return pbComponents
}
