/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Role_Type } from '@/client/graphql/generated';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import HomeIcon from '@mui/icons-material/Home';
import LocationIcon from '@mui/icons-material/LocationOn';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import PermDataSettingIcon from '@mui/icons-material/PermDataSetting';
import RouterIcon from '@mui/icons-material/Router';
import SimCardIcon from '@mui/icons-material/SimCard';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';

export const NavList = [
  {
    name: 'Home',
    path: '/console/home',
    icon: HomeIcon,
    isPrivate: true,
    isFullScreen: false,
    forRoles: [
      Role_Type.RoleOwner,
      Role_Type.RoleAdmin,
      Role_Type.RoleVendor,
      Role_Type.RoleNetworkOwner,
    ],
  },
  {
    name: 'Sites',
    path: '/console/sites',
    icon: LocationIcon,
    isPrivate: true,
    isFullScreen: false,
    forRoles: [Role_Type.RoleOwner, Role_Type.RoleAdmin],
  },
  {
    name: 'Nodes',
    path: '/console/nodes',
    icon: RouterIcon,
    isPrivate: true,
    isFullScreen: false,
    forRoles: [Role_Type.RoleOwner, Role_Type.RoleAdmin],
  },
  {
    name: 'Subscribers',
    path: '/console/subscribers',
    icon: SubscriberIcon,
    isPrivate: true,
    isFullScreen: false,
    forRoles: [
      Role_Type.RoleOwner,
      Role_Type.RoleAdmin,
      Role_Type.RoleVendor,
      Role_Type.RoleNetworkOwner,
    ],
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
    name: 'Members',
    path: '/manage/members',
    icon: SubscriberIcon,
  },
  {
    id: 'manage-sim',
    name: 'SIM pool',
    path: '/manage/sims',
    icon: SimCardIcon,
  },
  {
    id: 'manage-node',
    name: 'Node pool',
    path: '/manage/nodes',
    icon: AccountTreeIcon,
  },
  {
    id: 'manage-data-plan',
    name: 'Data plans',
    path: '/manage/data-plans',
    icon: PermDataSettingIcon,
  },
];

export const MANAGE_MENU_LIST_SMALL = [
  {
    id: 'manage-members',
    name: 'Members',
    path: '/manage/members',
    icon: SubscriberIcon,
  },
  {
    id: 'manage-sim',
    name: 'SIM pool',
    path: '/manage/sims',
    icon: SimCardIcon,
  },
  {
    id: 'manage-node',
    name: 'Node pool',
    path: '/manage/nodes',
    icon: AccountTreeIcon,
  },
  {
    id: 'manage-data-plan',
    name: 'Data plans',
    path: '/manage/data-plans',
    icon: PermDataSettingIcon,
  },
];
export const MY_ACCOUNT_MENU_LIST = [
  {
    id: 'billing',
    name: 'Billing',
    path: '/manage/billing',
    icon: AttachMoneyIcon,
  },
];
