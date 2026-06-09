/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Nodes — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const NodeSchema = z.object({
  id: z.string(),
  serial: z.string(),
  /** Human-friendly node name from the registry; undefined in mocks. */
  name: z.string().optional(),
  /** Raw connectivity (Online/Offline/Unknown) for the status dot. */
  connectivity: z.string().optional(),
  /** Raw operational state (Operational/Configured/Faulty/Unknown). */
  state: z.string().optional(),
  type: z.enum(['Tower node', 'Amplifier node', 'Controller node', 'Home node']),
  site: z.string(),
  status: z.enum(['online', 'degraded', 'offline', 'configuring']),
  cpu: z.number(),
  mem: z.number(),
  temp: z.number().nullable(),
  fw: z.string(),
  up: z.string(),
  note: z.string().optional(),
});

export type UkamaNode = z.infer<typeof NodeSchema>;

export const NODES: UkamaNode[] = [
  { id: 'n1', serial: 'uk-tnode-a05-1101', type: 'Tower node', site: 'Lusaka North', status: 'online', cpu: 34, mem: 51, temp: 42, fw: '13.2.1', up: '62d' },
  { id: 'n2', serial: 'uk-anode-a05-1102', type: 'Amplifier node', site: 'Lusaka North', status: 'online', cpu: 18, mem: 39, temp: 38, fw: '13.2.1', up: '62d' },
  { id: 'n3', serial: 'uk-tnode-a05-1140', type: 'Tower node', site: 'Kafue Bridge', status: 'degraded', cpu: 71, mem: 80, temp: 61, fw: '13.1.0', up: '9d', note: 'High latency' },
  { id: 'n4', serial: 'uk-tnode-a05-1166', type: 'Tower node', site: 'Chongwe East', status: 'offline', cpu: 0, mem: 0, temp: null, fw: '13.1.0', up: '—' },
  { id: 'n5', serial: 'uk-anode-a05-1167', type: 'Amplifier node', site: 'Chongwe East', status: 'offline', cpu: 0, mem: 0, temp: null, fw: '13.1.0', up: '—' },
  { id: 'n6', serial: 'uk-tnode-a05-1180', type: 'Tower node', site: 'Mumbwa Hub', status: 'online', cpu: 29, mem: 44, temp: 40, fw: '13.2.1', up: '40d' },
  { id: 'n7', serial: 'uk-tnode-a05-1210', type: 'Tower node', site: 'Kabwe Central', status: 'online', cpu: 41, mem: 55, temp: 45, fw: '13.2.1', up: '88d' },
  { id: 'n8', serial: 'uk-tnode-a05-1233', type: 'Tower node', site: 'Chilanga South', status: 'configuring', cpu: 5, mem: 22, temp: 36, fw: '13.2.1', up: '2h' },
];
