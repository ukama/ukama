/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** MUI module augmentation — custom variants (BUILD-PLAN §7.2 C). */
import '@mui/material/Chip';

declare module '@mui/material/Chip' {
  interface ChipPropsVariantOverrides {
    /** status badge tones (one status vocabulary, design finding #4) */
    ok: true;
    warn: true;
    err: true;
    neut: true;
    info: true;
    /** filter chip row (status filters, date chip) */
    chipFilter: true;
  }
}
