/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Business-lens datasets — ported from the prototype (`biz-data.jsx`,
 * modelled on the Maiko network screenshots) plus BUSINESS/KPIS (`data.jsx`).
 */
import type { KpiProps } from '@/components/Kpi';
import { series } from '@/lib/series';

export const KPIS = {
  health: 86,
  subs: 1284,
  subsSeries: series(1180, 14, 0.02, 0.09),
  data: 312,
  dataUnit: 'GB',
  dataSeries: series(280, 14, 0.1, 0.12),
  sales: 4820,
  salesSeries: series(4100, 14, 0.05, 0.18),
  uptime: 98.4,
  uptimeSeries: series(98.6, 14, 0.004, -0.002),
};

export interface BizSite {
  id: string;
  name: string;
  status: 'online' | 'warning' | 'offline';
  revenue: number;
  revToday: number;
  customers: number;
  custToday: number;
  data: string;
  uptime: number;
  top: string;
  issue: string | null;
  lat: number;
  lng: number;
}

export const BIZ_SITES: BizSite[] = [
  { id: 'clinicA', name: 'Clinic A', status: 'online', revenue: 920, revToday: 42, customers: 83, custToday: 83, data: '410 GB', uptime: 99.1, top: '1GB / 7d', issue: null, lat: -4.32, lng: 15.31 },
  { id: 'marketB', name: 'Market B', status: 'warning', revenue: 740, revToday: 31, customers: 91, custToday: 41, data: '360 GB', uptime: 94.8, top: '5GB / 30d', issue: 'Backhaul', lat: -2.52, lng: 23.62 },
  { id: 'schoolC', name: 'School C', status: 'online', revenue: 410, revToday: 18, customers: 42, custToday: 22, data: '190 GB', uptime: 99.8, top: '1GB / 7d', issue: null, lat: -1.68, lng: 27.48 },
  { id: 'clinicD', name: 'Clinic D', status: 'offline', revenue: 0, revToday: 0, customers: 83, custToday: 83, data: '0 GB', uptime: 82.2, top: '—', issue: 'Node down', lat: -6.12, lng: 25.21 },
  { id: 'villageE', name: 'Village E', status: 'online', revenue: 280, revToday: 12, customers: 17, custToday: 17, data: '90 GB', uptime: 98.9, top: 'Night', issue: null, lat: -3.02, lng: 19.18 },
];

export const BIZ_HOME = {
  kpis: [
    { label: 'Revenue', value: '$126.40', delta: '18% vs yesterday', dir: 'up' },
    { label: 'Active customers', value: '142', delta: '+18 today', dir: 'up' },
    { label: 'Data sold', value: '87 GB', sub: '1.8 TB this month' },
    { label: 'Network uptime', value: '96.8%', sub: '6/7 sites online' },
  ] as KpiProps[],
  topPackages: [
    { name: '1GB / 7 days', revenue: 840, sold: 420, color: 'var(--uk-ac)' },
    { name: '5GB / 30 days', revenue: 760, sold: 95, color: 'var(--uk-secondary)' },
    { name: 'Night bundle', revenue: 175, sold: 175, color: 'var(--uk-success-bright)' },
  ],
};

export const BIZ_SALES = {
  kpis: [
    { label: 'Revenue this month', value: '$3,420', delta: '12% vs last month', dir: 'up' },
    { label: 'Purchases', value: '812', delta: '+64 this week', dir: 'up' },
    { label: 'Avg purchase', value: '$4.21', sub: 'stable' },
    { label: 'Paid customers', value: '376', delta: '+28 this month', dir: 'up' },
  ] as KpiProps[],
  trend: [1820, 2080, 1940, 2460, 2520, 2700, 2980, 3180, 3420],
  bySite: [
    { name: 'Clinic A', value: 920, color: 'var(--uk-ac)' },
    { name: 'Market B', value: 740, color: 'var(--uk-secondary)' },
    { name: 'School C', value: 410, color: 'var(--uk-success-bright)' },
    { name: 'Village E', value: 280, color: 'var(--uk-orange)' },
  ],
  byPackage: [
    { name: '1GB / 7 days', value: 1240, color: 'var(--uk-ac)' },
    { name: '5GB / 30 days', value: 880, color: 'var(--uk-secondary)' },
    { name: 'Night bundle', value: 390, color: 'var(--uk-success-bright)' },
    { name: '2GB weekend', value: 220, color: 'var(--uk-orange)' },
  ],
};

export interface BizPackageRow {
  pkg: string;
  price: string;
  validity: string;
  sold: number;
  revenue: string;
  data: string;
  status: 'active' | 'testing' | 'lowsales';
}

export const BIZ_PACKAGES = {
  rows: [
    { pkg: '1GB / 7 days', price: '$2', validity: '7 days', sold: 420, revenue: '$840', data: '390 GB', status: 'active' },
    { pkg: '5GB / 30 days', price: '$8', validity: '30 days', sold: 95, revenue: '$760', data: '430 GB', status: 'active' },
    { pkg: 'Night bundle', price: '$1', validity: '12 hours', sold: 175, revenue: '$175', data: '210 GB', status: 'active' },
    { pkg: '2GB weekend', price: '$3', validity: '2 days', sold: 72, revenue: '$216', data: '180 GB', status: 'testing' },
    { pkg: '10GB / 30 days', price: '$14', validity: '30 days', sold: 18, revenue: '$252', data: '170 GB', status: 'lowsales' },
  ] as BizPackageRow[],
  mix: [
    { name: '1GB / 7 days', value: 840, color: 'var(--uk-ac)' },
    { name: '5GB / 30 days', value: 760, color: 'var(--uk-secondary)' },
    { name: 'Night', value: 175, color: 'var(--uk-success-bright)' },
  ],
};

export interface BizResourceRow {
  res: string;
  type: string;
  status: 'online' | 'warning' | 'offline';
  site: string;
  affected: number;
  context: string;
  updated: string;
}

export const BIZ_NETWORK = {
  kpis: [
    { label: 'Network uptime', value: '96.8%', sub: 'today' },
    { label: 'Sites online', value: '6/7', sub: 'one offline', danger: true },
    { label: 'Nodes online', value: '12/13', sub: 'one down', danger: true },
    { label: 'Customers affected', value: '83', sub: 'Clinic D', danger: true },
  ] as KpiProps[],
  rows: [
    { res: 'Node T-001', type: 'Tower', status: 'offline', site: 'Clinic D', affected: 83, context: '$42 today', updated: '14 min ago' },
    { res: 'Backhaul B-2', type: 'Backhaul', status: 'warning', site: 'Market B', affected: 41, context: 'Speeds slow', updated: '7 min ago' },
    { res: 'Power P-4', type: 'Power', status: 'warning', site: 'Clinic A', affected: 83, context: 'Risk later', updated: '5 min ago' },
    { res: 'Node A-006', type: 'Amplifier', status: 'online', site: 'School C', affected: 0, context: 'Healthy', updated: '1 min ago' },
  ] as BizResourceRow[],
  summary: [
    { tone: 'err', title: 'Offline node', detail: 'Clinic D · 83 customers affected' },
    { tone: 'warn', title: 'Backhaul warning', detail: 'Market B · possible slow speeds' },
    { tone: 'warn', title: 'Power warning', detail: 'Clinic A · battery low' },
  ] as { tone: 'err' | 'warn' | 'ok' | 'info'; title: string; detail: string }[],
};

export const BIZ_INVENTORY = {
  kpis: [
    { label: 'Available SIMs', value: '120', sub: 'ready to assign' },
    { label: 'Active SIMs', value: '376', sub: 'in use' },
    { label: 'Available nodes', value: '3', sub: 'in inventory' },
    { label: 'Deployed nodes', value: '12', sub: 'across sites' },
  ] as KpiProps[],
  sims: [
    { iccid: '8923…0089', status: 'available', cust: '—', site: '—', date: '—', issue: '—' },
    { iccid: '8923…0021', status: 'active', cust: '+243…908', site: 'Clinic A', date: 'May 21', issue: '—' },
    { iccid: '8923…0044', status: 'active', cust: '+243…441', site: 'Market B', date: 'May 22', issue: '—' },
    { iccid: '8923…0091', status: 'failed', cust: '—', site: '—', date: '—', issue: 'Activation failed' },
    { iccid: '8923…0124', status: 'assigned', cust: '+243…128', site: 'Clinic A', date: 'Pending', issue: 'Waiting' },
  ],
  nodes: [
    { serial: 'Node A-006', type: 'Amplifier', status: 'deployed', site: 'School C', date: 'May 12' },
    { serial: 'Node T-009', type: 'Tower', status: 'deployed', site: 'Village E', date: 'May 09' },
    { serial: 'Node T-014', type: 'Tower', status: 'available', site: '—', date: 'May 24' },
    { serial: 'Node A-018', type: 'Amplifier', status: 'available', site: '—', date: 'May 24' },
  ],
};

export const BIZ_SITE_DETAIL = {
  meta: 'Nodes: 2/2 online · Top package: 1GB / 7 days',
  kpis: [
    { label: 'Revenue this month', value: '$920', delta: '9%', dir: 'up' },
    { label: 'Active customers', value: '83', delta: '+8 today', dir: 'up' },
    { label: 'Data sold', value: '410 GB', sub: 'this month' },
    { label: 'Uptime', value: '99.1%', sub: 'last 30 days' },
  ] as { label: string; value: string; delta?: string; dir?: 'up' | 'down'; sub?: string }[],
  tabs: ['Overview', 'Nodes', 'Subscribers', 'Metrics', 'Config'],
  resources: [
    { res: 'Tower node', status: 'online', seen: '2 min ago', issue: '—' },
    { res: 'Amplifier', status: 'online', seen: '2 min ago', issue: '—' },
    { res: 'Backhaul', status: 'online', seen: '1 min ago', issue: '—' },
    { res: 'Power', status: 'warning', seen: '5 min ago', issue: 'Battery low' },
  ],
};

export interface AdminItem {
  name: string;
  desc: string;
  href: string;
}

export const BIZ_ADMIN: { group: string; items: AdminItem[] }[] = [
  {
    group: 'Infrastructure',
    items: [
      { name: 'Nodes', desc: 'Manage node inventory, status, assignment', href: '/business/inventory' },
      { name: 'Sites admin', desc: 'Configure site deployment settings', href: '/business/sites' },
      { name: 'Networks', desc: 'Manage network settings', href: '/business/settings' },
      { name: 'SIMs', desc: 'SIM inventory and allocation', href: '/business/manage/sim-pool' },
    ],
  },
  {
    group: 'People',
    items: [
      { name: 'Members', desc: 'Users, roles, access', href: '/business/manage/members' },
      { name: 'Roles', desc: 'Permission management', href: '/business/manage/members' },
    ],
  },
  {
    group: 'Billing',
    items: [
      { name: 'Billing', desc: 'Payment method and account settings', href: '/business/manage/billing' },
      { name: 'Invoices', desc: 'Billing history', href: '/business/manage/billing' },
    ],
  },
  {
    group: 'Configuration',
    items: [
      { name: 'Organization', desc: 'Company / account settings', href: '/business/settings' },
      { name: 'Network config', desc: 'Defaults and advanced settings', href: '/business/settings' },
    ],
  },
];
