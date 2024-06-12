/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import HomeIcon from '@mui/icons-material/Home';
import LocationIcon from '@mui/icons-material/LocationOn';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import RouterIcon from '@mui/icons-material/Router';

export const NavList = [
  {
    name: 'Home',
    path: '/home',
    icon: HomeIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Sites',
    path: '/sites',
    icon: LocationIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Nodes',
    path: '/nodes',
    icon: RouterIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Subscribers',
    path: '/subscribers',
    icon: SubscriberIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  // {
  //   name: 'Site Planning',
  //   path: '/site_planning',
  //   icon: LayersIcon,
  //   isPrivate: true,
  //   isFullScreen: false,
  // },
];
