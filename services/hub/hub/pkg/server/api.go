package server

import (
	"github.com/ukama/ukama/services/hub/hub/pkg"
	"time"
)

type CAppRequest struct {
	Name string `path:"name" validate:"required"`
	// version + extension
	ArtifactName string `path:"filename" validate:"required"`
}

type VersionListRequest struct {
	Name string `path:"name" validate:"required"`
}

type VersionListResponse struct {
	Name     string         `json:"name"`
	Versions *[]VersionInfo `json:"artifacts"`
}

type VersionInfo struct {
	Version string       `json:"version"`
	Formats []FormatInfo `json:"formats"`
}

type FormatInfo struct {
	Type      string            `json:"type"`
	Url       string            `json:"url"`
	CreatedAt time.Time         `json:"created_at"`
	SizeBytes int64             `json:"size_bytes,omitempty"`
	ExtraInfo map[string]string `json:"extra_info,omitempty"`
}

type CAppArtifact struct {
	Name string `path:"name" validate:"required"`
	// version + extension
	ArtifactName string `path:"filename" validate:"required"`
}

type CAppsListResponse struct {
	Artifacts *[]pkg.CappInfo `json:"capps"`
}

type CApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
