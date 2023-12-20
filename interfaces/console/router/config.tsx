/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { CenterContainer } from '@/styles/global';
import { CircularProgress } from '@mui/material';
import { ComponentType, LazyExoticComponent, ReactNode, lazy } from 'react';

import HomeIcon from '@mui/icons-material/Home';
import LocationIcon from '@mui/icons-material/LocationOn';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import RouterIcon from '@mui/icons-material/Router';

export interface IRoute {
  // Path, like in basic prop
  path: string;
  // Exact, like in basic prop
  exact: boolean;
  // Preloader for lazy loading
  fallback: NonNullable<ReactNode> | null;
  // Lazy Loaded component
  component?: LazyExoticComponent<ComponentType<any>>;
  // Sub routes
  routes?: IRoute[];
  // Redirect path
  redirect?: string;
  // If router is private, this is going to be true
  private?: boolean;

  isFullScreen?: boolean;
}

const Loader = (
  <CenterContainer>
    <CircularProgress />
  </CenterContainer>
);

const getRouteObject = (
  path = '/',
  component = 'OnBoarding',
  isPrivate = true,
  isFullScreen = false,
) => {
  return {
    path: path,
    exact: true,
    fallback: Loader,
    private: isPrivate,
    isFullScreen: isFullScreen,
    component: lazy(() => import(`../pages/${component}`)),
  };
};

export const routes = {
  Root: getRouteObject('/home', 'Home', true),
  Nodes: getRouteObject('/nodes', 'Nodes', true),
  Sites: getRouteObject('/sites', 'Sites', true),
  Users: getRouteObject('/users', 'Users', true),
  Settings: getRouteObject('/settings', 'Settings', true, true),
  Billing: getRouteObject('/billing', 'Billing', true),
  OnBoarding: getRouteObject('/', 'OnBoarding', true, true),
  //Public Routes//
  Error: getRouteObject('/*', 'ErrorPage', true, true),
  //
};

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
