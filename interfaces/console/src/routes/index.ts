/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import AccountTreeIcon from '@mui/icons-material/AccountTree';
import HomeIcon from '@mui/icons-material/Home';
import LocationIcon from '@mui/icons-material/LocationOn';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import PermDataSettingIcon from '@mui/icons-material/PermDataSetting';
import RouterIcon from '@mui/icons-material/Router';
import SimCardIcon from '@mui/icons-material/SimCard';

export const NavList = [
  {
    name: 'Home',
    path: '/console/home',
    icon: HomeIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Sites',
    path: '/console/sites',
    icon: LocationIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Nodes',
    path: '/console/nodes',
    icon: RouterIcon,
    isPrivate: true,
    isFullScreen: false,
  },
  {
    name: 'Subscribers',
    path: '/console/subscribers',
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

export const MANAGE_MENU_LIST = [
  {
    id: 'manage-members',
    name: 'Manage members',
    path: '/manage/members',
    icon: SubscriberIcon,
  },
  {
    id: 'manage-sim',
    name: 'Manage sim pool',
    path: '/manage/sims',
    icon: SimCardIcon,
  },
  {
    id: 'manage-node',
    name: 'Manage node pool',
    path: '/manage/nodes',
    icon: AccountTreeIcon,
  },
  {
    id: 'manage-data-plan',
    name: 'Manage data plans',
    path: '/manage/data-plans',
    icon: PermDataSettingIcon,
  },
];
