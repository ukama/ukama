/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Sites — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const SiteSchema = z.object({
  id: z.string(),
  name: z.string(),
  area: z.string(),
  status: z.enum(['online', 'degraded', 'offline']),
  subs: z.number(),
  nodes: z.number(),
  uptime: z.number(),
  battery: z.number(),
  signal: z.number().nullable(),
  data: z.string(),
  /** dummy coordinates within the operating region (real geo at API phase) */
  lat: z.number(),
  lng: z.number(),
  plan: z.string(),
  issue: z.string().optional(),
});

export type Site = z.infer<typeof SiteSchema>;

export const SITES: Site[] = [
  {
    id: 's1',
    name: 'Lusaka North',
    area: 'Matero District',
    status: 'online',
    subs: 412,
    nodes: 4,
    uptime: 99.97,
    battery: 96,
    signal: -71,
    data: '1.8 TB',
    lat: -15.32,
    lng: 28.35,
    plan: 'Grid + solar',
  },
  {
    id: 's2',
    name: 'Kafue Bridge',
    area: 'Kafue Road',
    status: 'degraded',
    subs: 188,
    nodes: 3,
    uptime: 97.2,
    battery: 41,
    signal: -92,
    data: '612 GB',
    lat: -15.98,
    lng: 27.92,
    plan: 'Solar',
    issue: 'Backhaul latency high · battery 41%',
  },
  {
    id: 's3',
    name: 'Chongwe East',
    area: 'Great East Road',
    status: 'offline',
    subs: 74,
    nodes: 2,
    uptime: 62.4,
    battery: 8,
    signal: null,
    data: '0 GB',
    lat: -15.28,
    lng: 29.05,
    plan: 'Solar',
    issue: 'Offline 2h 14m · battery critical',
  },
  {
    id: 's4',
    name: 'Mumbwa Hub',
    area: 'Mumbwa Town',
    status: 'online',
    subs: 233,
    nodes: 3,
    uptime: 99.81,
    battery: 88,
    signal: -76,
    data: '940 GB',
    lat: -14.98,
    lng: 27.06,
    plan: 'Grid',
  },
  {
    id: 's5',
    name: 'Kabwe Central',
    area: 'Freedom Way',
    status: 'online',
    subs: 301,
    nodes: 4,
    uptime: 99.93,
    battery: 91,
    signal: -69,
    data: '1.2 TB',
    lat: -14.44,
    lng: 28.45,
    plan: 'Grid + solar',
  },
  {
    id: 's6',
    name: 'Chilanga South',
    area: 'Kafue Road S',
    status: 'degraded',
    subs: 96,
    nodes: 2,
    uptime: 98.1,
    battery: 67,
    signal: -88,
    data: '318 GB',
    lat: -15.62,
    lng: 28.22,
    plan: 'Solar',
    issue: '1 of 2 nodes configuring',
  },
];
