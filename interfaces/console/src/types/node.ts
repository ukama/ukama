/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

export type TNodeSiteTree = {
  id: string;
  name: string;
  nodeId: string;
  nodeType: string;
  nodeName: string;
};

export type TNodePoolData = {
  id: string;
  type: string;
  site: string;
  state: string;
  network: string;
  createdAt: string;
  connectivity: string;
};
