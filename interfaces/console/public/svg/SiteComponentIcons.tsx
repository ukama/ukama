/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';

interface IconProps {
  size?: number;
  color?: string;
}

export const NodeIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="2" y="6" width="20" height="12" rx="2"></rect>
    <line x1="6" y1="12" x2="10" y2="12"></line>
    <line x1="14" y1="12" x2="18" y2="12"></line>
  </svg>
);

export const SwitchIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="2" y="6" width="20" height="12" rx="2"></rect>
    <line x1="6" y1="10" x2="6" y2="14"></line>
    <line x1="10" y1="10" x2="10" y2="14"></line>
    <line x1="14" y1="10" x2="14" y2="14"></line>
    <line x1="18" y1="10" x2="18" y2="14"></line>
  </svg>
);

export const BackhaulIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M5 12.55a11 11 0 0 1 14.08 0"></path>
    <path d="M1.42 9a16 16 0 0 1 21.16 0"></path>
    <path d="M8.53 16.11a6 6 0 0 1 6.95 0"></path>
    <circle cx="12" cy="20" r="2"></circle>
  </svg>
);

export const ChargeControllerIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="2" y="4" width="20" height="16" rx="2"></rect>
    <circle cx="12" cy="12" r="4"></circle>
    <line x1="12" y1="8" x2="12" y2="16"></line>
    <line x1="8" y1="12" x2="16" y2="12"></line>
  </svg>
);

export const SolarPanelIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="3" y="4" width="18" height="12" rx="2"></rect>
    <line x1="3" y1="8" x2="21" y2="8"></line>
    <line x1="3" y1="12" x2="21" y2="12"></line>
    <line x1="9" y1="4" x2="9" y2="16"></line>
    <line x1="15" y1="4" x2="15" y2="16"></line>
    <line x1="12" y1="16" x2="12" y2="20"></line>
    <line x1="8" y1="20" x2="16" y2="20"></line>
  </svg>
);

export const BatteryIcon: React.FC<IconProps> = ({
  size = 24,
  color = 'currentColor',
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke={color}
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <rect x="2" y="7" width="18" height="10" rx="2"></rect>
    <line x1="22" y1="11" x2="22" y2="13"></line>
    <line x1="6" y1="10" x2="6" y2="14"></line>
    <line x1="10" y1="10" x2="10" y2="14"></line>
  </svg>
);
