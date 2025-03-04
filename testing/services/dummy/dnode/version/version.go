/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package version

// Minor version is autoupdated by the build system
// NOTE: use go build -ldflags "-X github.com/ukama/ukama/systems/registry/network/cmd/version.Version==$(git describe)".
var Version = "v0.0.debug"
