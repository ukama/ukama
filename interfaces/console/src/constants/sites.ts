/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { NODE_ACTIONS_ENUM } from './nodes';

export const INSTALLATION_FLOW = 'ins';
export const CHECK_SITE_FLOW = 'chk';

export const SITE_ACTIONS_BUTTONS = [
  {
    id: NODE_ACTIONS_ENUM.TOGGLE_RADIO,
    name: 'Toggle Radio',
    type: 'toggle',
    consent: 'Are you sure you want to toggle the radio?',
  },
  {
    id: NODE_ACTIONS_ENUM.TOGGLE_SERVICE,
    name: 'Toggle Service',
    type: 'toggle',
    consent: 'Are you sure you want to toggle the service?',
  },
];

export const SITE_CONFIG_STEPS = [
  'Configure Site Installation (1/2)',
  'Configure Site Installation (2/2)',
];

export const SITE_STATUS = {
  ONLINE: 'Online',
  OFFLINE: 'Offline',
  WARNING: 'Warning',
};

export const TIME_FILTER_OPTIONS = [
  { id: '1', label: 'LIVE' },
  { id: '2', label: 'ZOOM' },
];

export const SITE_KPI_TYPES = {
  SITE_UPTIME: 'site_uptime_seconds',
  BATTERY_CHARGE_PERCENTAGE: 'battery_charge',
  BACKHAUL_SPEED: 'backhaul_downlink',
  SITE_UPTIME_PERCENTAGE: 'site_uptime_percentage',
  NODE_UPTIME: 'uptime',
  ACTIVE_SUBSCRIBERS: 'node_active_subscribers',
};

export interface SiteKpiConfig {
  id: string;
  name: string;
  unit: string;
  description: string;
  port?: number;
  tickInterval?: number;
  tickPositions?: number[];
  threshold?: {
    min: number;
    normal: number;
    max: number;
  } | null;
  show?: boolean;
  format: string;
}

export interface SectionData {
  [key: string]: SiteKpiConfig[];
}

export const SITE_KPIS = {
  SOLAR: {
    metrics: [
      {
        unit: 'W',
        show: true,
        name: 'Solar panel power',
        id: 'solar_panel_power',
        description: 'Solar power generation',
        tickInterval: 100,
        min: 0,
        max: 600,
        threshold: {
          min: 100,
          normal: 300,
          max: 500,
        },
        format: 'number',
        tickPositions: [0, 100, 200, 300, 400, 500, 600],
      },
      {
        unit: 'V',
        show: true,
        name: 'Solar panel voltage',
        id: 'solar_panel_voltage',
        description: 'Solar panel voltage',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 50,
          normal: 75,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
      {
        unit: 'A',
        show: true,
        name: 'Solar panel current',
        id: 'solar_panel_current',
        description: 'Solar panel current',
        tickInterval: 2,
        min: 0,
        max: 12,
        threshold: {
          min: 2,
          normal: 5,
          max: 8,
        },
        format: 'number',
        tickPositions: [0, 2, 4, 6, 8, 10, 12],
      },
    ],
  },
  BATTERY: {
    metrics: [
      {
        unit: '%',
        show: true,
        name: 'Battery charge',
        id: 'battery_charge',
        description: 'Battery charge percentage',
        tickInterval: 20,
        min: 0,
        max: 120,
        threshold: {
          min: 20,
          normal: 70,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  CONTROLLER: {
    metrics: [
      {
        unit: 'W',
        show: true,
        name: 'Controller temperature',
        id: 'controller_temperature',
        description: 'Controller temperature',
        tickInterval: 20,
        min: 0,
        max: 80,
        threshold: {
          min: 0,
          normal: 40,
          max: 80,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80],
      },
      {
        unit: 'V',
        show: true,
        name: 'Controller load current',
        id: 'load_current',
        description: 'Controller load current',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 50,
          normal: 75,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  MAIN_BACKHAUL: {
    metrics: [
      {
        unit: 'ms',
        show: true,
        name: 'Backhaul latency',
        id: 'backhaul_latency',
        description: 'Backhaul latency',
        tickInterval: 500,
        min: 0,
        max: 3000,
        threshold: {
          min: 0,
          normal: 500,
          max: 1000,
        },
        format: 'number',
        lowerIsBetter: true,
        tickPositions: [0, 500, 1000, 1500, 2000, 2500, 3000],
      },
      {
        unit: 'Mbps',
        show: true,
        name: 'Backhaul downlink',
        id: 'backhaul_downlink',
        description: 'Backhaul downlink goodput',
        tickInterval: 10,
        min: 0,
        max: 120,
        threshold: {
          min: 20,
          normal: 50,
          max: 100,
        },
        format: 'number',
        tickPositions: [0, 20, 40, 60, 80, 100, 120],
      },
    ],
  },
  SWITCH: {
    metrics: [
      {
        unit: 'bps',
        show: true,
        name: 'Tower node switch port speed',
        id: 'switch_port_1_speed',
        description: 'Tower node switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000000000,
        threshold: {
          min: 0,
          normal: 500000000,
          max: 1000000000,
        },
        format: 'number',
        port: 1,
      },
      {
        unit: 'W',
        show: true,
        name: 'Tower node switch port power',
        id: 'switch_port_1_power',
        description: 'Tower node switch port power',
        tickInterval: 1,
        min: 0,
        max: 24,
        threshold: {
          min: 0,
          normal: 12,
          max: 24,
        },
        format: 'number',
        port: 1,
      },
      {
        unit: 'bps',
        show: true,
        name: 'Amplifier node switch port speed',
        id: 'switch_port_2_speed',
        description: 'Amplifier node switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000000000,
        threshold: {
          min: 0,
          normal: 500000000,
          max: 1000000000,
        },
        format: 'number',
        port: 2,
      },
      {
        unit: 'W',
        show: true,
        name: 'Amplifier node switch port power',
        id: 'switch_port_2_power',
        description: 'Amplifier node switch port power',
        tickInterval: 1,
        min: 0,
        max: 24,
        threshold: {
          min: 0,
          normal: 12,
          max: 24,
        },
        format: 'number',
        port: 2,
      },
      {
        unit: 'bps',
        show: true,
        name: 'Controller node switch port speed',
        id: 'switch_port_3_speed',
        description: 'Controller node switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000000000,
        threshold: {
          min: 0,
          normal: 500000000,
          max: 1000000000,
        },
        format: 'number',
        port: 3,
      },
      {
        unit: 'W',
        show: true,
        name: 'Controller node switch port power',
        id: 'switch_port_3_power',
        description: 'Controller node switch port power',
        tickInterval: 1,
        min: 0,
        max: 24,
        threshold: {
          min: 0,
          normal: 12,
          max: 24,
        },
        format: 'number',
        port: 3,
      },
      {
        unit: 'bps',
        show: true,
        name: 'Backhaul node switch port speed',
        id: 'switch_port_9_speed',
        description: 'Backhaul node switch port speed',
        tickInterval: 100,
        min: 0,
        max: 1000000000,
        threshold: {
          min: 0,
          normal: 500000000,
          max: 1000000000,
        },
        format: 'number',
        port: 9,
      },
      {
        unit: 'W',
        show: true,
        name: 'Backhaul node switch port power',
        id: 'switch_port_9_power',
        description: 'Backhaul node switch port power',
        tickInterval: 1,
        min: 0,
        max: 24,
        threshold: {
          min: 0,
          normal: 12,
          max: 24,
        },
        format: 'number',
        port: 9,
      },
    ],
  },
  SITE: {
    stats: [
      {
        unit: '',
        show: true,
        name: 'Active subscribers',
        format: 'number',
        id: 'node_active_subscribers',
        description: 'Current active subscribers on the site',
        threshold: null,
      },
    ],
  },
};
