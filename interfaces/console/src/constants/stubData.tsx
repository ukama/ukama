/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
export const Dummysites = [
  {
    siteId: 'site1',
    name: 'Site very long name #1',
    address: '1000 Nelson Way',
    users: 2,
    status: {
      online: true,
      charging: true,
      signal: 'Good',
    },
  },
  {
    siteId: 'site2',
    name: 'Site very long name #2',
    address: '2000 Nelson Way',
    users: 3,
    status: {
      online: false,
      charging: true,
      signal: 'Poor',
    },
  },
  {
    siteId: 'site3',
    name: 'Site very long name #3',
    address: '3000 Nelson Way',
    users: 1,
    status: {
      online: true,
      charging: false,
      signal: 'Good',
    },
  },
];
export const NodeApps = [
  {
    id: 1,
    nodeAppName: 'App # 1',
    cpu: 12,
    memory: 30,
    version: '01',
  },
  {
    id: 2,
    nodeAppName: 'App # 2',
    cpu: 12,
    memory: 30,
    version: '01',
  },
  {
    id: 3,
    nodeAppName: 'App # 3',
    cpu: 12,
    memory: 30,
    version: '01',
  },
  {
    id: 4,
    nodeAppName: 'App # 4',
    cpu: 30,
    memory: 30,
    version: '01',
  },
];
