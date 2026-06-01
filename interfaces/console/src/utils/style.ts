/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

export const hexToRGB = (hex: string, alpha: number): string => {
  const h = '0123456789ABCDEF';
  const r = h.indexOf(hex[1]) * 16 + h.indexOf(hex[2]);
  const g = h.indexOf(hex[3]) * 16 + h.indexOf(hex[4]);
  const b = h.indexOf(hex[5]) * 16 + h.indexOf(hex[6]);
  if (alpha) {
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  return `rgba(${r}, ${g}, ${b})`;
};
