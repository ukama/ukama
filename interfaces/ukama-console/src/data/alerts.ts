/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Attention feed / notifications — ported from the design prototype (`data.jsx`). */
import { z } from 'zod';

export const AlertSchema = z.object({
  id: z.string(),
  sev: z.enum(['critical', 'warning', 'info']),
  icon: z.string(),
  title: z.string(),
  detail: z.string(),
  site: z.string().optional(),
  action: z.string(),
  age: z.string(),
});

export type Alert = z.infer<typeof AlertSchema>;

export const ALERTS: Alert[] = [
  {
    id: 'a1',
    sev: 'critical',
    icon: 'cell_tower',
    title: 'Chongwe East is offline',
    detail: 'No heartbeat for 2h 14m · battery critical (8%)',
    site: 'Chongwe East',
    action: 'Diagnose',
    age: '2h',
  },
  {
    id: 'a2',
    sev: 'warning',
    icon: 'battery_alert',
    title: 'Kafue Bridge battery low',
    detail: '41% and falling · solar charge interrupted',
    site: 'Kafue Bridge',
    action: 'View site',
    age: '38m',
  },
  {
    id: 'a3',
    sev: 'warning',
    icon: 'network_check',
    title: 'Backhaul latency high',
    detail: 'uk-tnode-…1140 averaging 210 ms over 1h',
    site: 'Kafue Bridge',
    action: 'View node',
    age: '1h',
  },
  {
    id: 'a4',
    sev: 'info',
    icon: 'sync',
    title: 'Node configuring at Chilanga South',
    detail: 'uk-tnode-…1233 provisioning · ~4 min remaining',
    site: 'Chilanga South',
    action: 'View node',
    age: 'just now',
  },
  {
    id: 'a5',
    sev: 'info',
    icon: 'sim_card',
    title: 'SIM pool running low',
    detail: '611 SIMs available · below 700 threshold',
    action: 'Order SIMs',
    age: '5h',
  },
];
