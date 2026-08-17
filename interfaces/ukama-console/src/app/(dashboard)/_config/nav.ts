/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Per-lens navigation config (BUILD-PLAN §2, prototype shell.jsx).
 * The money / no-money split: revenue lives in Business; running the
 * physical network lives in Network; the agent lens is deliberately tiny.
 */

export type Lens = 'business' | 'network' | 'customer';

export interface NavItem {
  href: string;
  label: string;
  /** icon key resolved via _components/icons.tsx */
  icon: string;
  /** exact match only (e.g. lens home) instead of prefix matching */
  exact?: boolean;
  /**
   * Opens outside the console (new tab, no active-state matching). The
   * Sidebar builds the real URL at render time — for the backend-status
   * item it appends ?org={orgName} to NEXT_PUBLIC_STATUS_APP_URL.
   */
  external?: boolean;
}

export interface NavGroup {
  group?: string;
  items: NavItem[];
}

export const LENSES: { id: Lens; label: string; icon: string; href: string }[] = [
  { id: 'business', label: 'Business', icon: 'insights', href: '/business' },
  { id: 'network', label: 'Network', icon: 'hub', href: '/network' },
  { id: 'customer', label: 'Customer', icon: 'badge', href: '/customer/customers' },
];

export const BIZ_NAV: NavGroup[] = [
  {
    items: [
      { href: '/business', label: 'Home', icon: 'home', exact: true },
      { href: '/business/revenue', label: 'Revenue', icon: 'payments' },
      { href: '/business/customers', label: 'Customers', icon: 'group' },
      { href: '/business/packages', label: 'Packages', icon: 'donut_small' },
    ],
  },
  {
    group: 'Manage',
    items: [
      { href: '/business/manage/data-plans', label: 'Data plans', icon: 'apps' },
      // Billing is hidden for now — will finish later. Re-enable this item
      // and remove the notFound() guard in business/manage/billing/page.tsx.
      // { href: '/business/manage/billing', label: 'Billing', icon: 'monetization_on' },
      { href: '/business/manage/members', label: 'Members', icon: 'manage_accounts' },
      { href: '/business/manage/sim-pool', label: 'SIM pool', icon: 'sim_card' },
    ],
  },
];

export const NETWORK_NAV: NavGroup[] = [
  {
    items: [
      { href: '/network', label: 'Home', icon: 'home', exact: true },
      { href: '/network/sites', label: 'Sites', icon: 'location_on' },
      { href: '/network/nodes', label: 'Nodes', icon: 'router' },
      { href: '/network/customers', label: 'Customers', icon: 'group' },
    ],
  },
  {
    group: 'Manage',
    items: [
      { href: '/network/manage/node-pool', label: 'Node pool', icon: 'account_tree' },
      { href: '/network/manage/sim-pool', label: 'SIM pool', icon: 'sim_card' },
    ],
  },
];

export const AGENT_NAV: NavGroup[] = [
  {
    items: [
      { href: '/customer/customers', label: 'Customers', icon: 'group' },
      { href: '/customer/data-plans', label: 'Data plans', icon: 'apps' },
    ],
  },
];

export const NAV_BY_LENS: Record<Lens, NavGroup[]> = {
  business: BIZ_NAV,
  network: NETWORK_NAV,
  customer: AGENT_NAV,
};

/** Lens from the first path segment; Business is the default lens. */
export function lensFromPath(pathname: string): Lens {
  if (pathname.startsWith('/network')) return 'network';
  if (pathname.startsWith('/customer')) return 'customer';
  return 'business';
}

/** Support is hidden in the agent lens (prototype shell.jsx). */
/** Sidebar sentinel: external link to the status app (URL built in Sidebar). */
export const BACKEND_STATUS_ITEM: NavItem = {
  href: '/backend-status',
  label: 'Service status',
  icon: 'monitor_heart',
  external: true,
};

export function bottomNav(lens: Lens): NavItem[] {
  const settings = { href: `/${lens}/settings`, label: 'Settings', icon: 'settings' };
  if (lens === 'customer') return [BACKEND_STATUS_ITEM, settings];
  return [
    { href: `/${lens}/support`, label: 'Support', icon: 'support_agent' },
    BACKEND_STATUS_ITEM,
    settings,
  ];
}
