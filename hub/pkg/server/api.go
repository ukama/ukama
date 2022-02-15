package server

import "github.com/ukama/ukamaX/hub/pkg"

type CAppRequest struct {
	Name    string `path:"name" validate:"required"`
	Version string `path:"version" validate:"required"`
}

type CAppListRequest struct {
	Name string `path:"name" validate:"required"`
}

type CAppListResponse struct {
	Artifacts *[]pkg.AritfactInfo `json:"artifacts"`
}

type CAppsListResponse struct {
	Artifacts *[]pkg.CappInfo `json:"capps"`
}

type CApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
