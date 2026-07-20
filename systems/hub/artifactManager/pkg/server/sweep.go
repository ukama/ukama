/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"time"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
)

// RunConsistencySweep periodically finds stored tar.gz artifacts that lack a chunk
// index (silent chunking failures) and re-chunks them. Runs until ctx is cancelled.
func (s *ArtifcatServer) RunConsistencySweep(ctx context.Context, interval time.Duration, types []string) {
	if interval <= 0 {
		interval = 10 * time.Minute
	}
	if len(types) == 0 {
		types = []string{"app", "cert"}
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Infof("Artifact consistency sweep running every %s for types %v", interval, types)
	for {
		select {
		case <-ctx.Done():
			log.Infof("Artifact consistency sweep stopped")
			return
		case <-ticker.C:
			s.sweepOnce(ctx, types)
		}
	}
}

func (s *ArtifcatServer) sweepOnce(ctx context.Context, types []string) {
	for _, aType := range types {
		apps, err := s.storage.ListApps(ctx, aType)
		if err != nil {
			log.Errorf("sweep: list apps for type %s failed: %v", aType, err)
			continue
		}
		for _, name := range apps {
			if name == "" {
				continue
			}
			versions, err := s.storage.ListVersions(ctx, name, aType)
			if err != nil {
				log.Errorf("sweep: list versions for %s/%s failed: %v", aType, name, err)
				continue
			}
			for i := range *versions {
				info := (*versions)[i]
				if info.Chunked {
					continue
				}
				v, err := semver.NewVersion(info.Version)
				if err != nil {
					continue // skip INVALID_VERSION_FORMAT
				}
				log.Warnf("sweep: %s/%s@%s has no chunk index; re-chunking", aType, name, v.String())
				if err := s.chunkAndIndex(ctx, name, aType, v); err != nil {
					log.Errorf("sweep: re-chunk %s/%s@%s failed: %v", aType, name, v.String(), err)
				}
			}
		}
	}
}
