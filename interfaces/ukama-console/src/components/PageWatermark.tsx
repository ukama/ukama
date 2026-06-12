/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Large, faint icon anchored to the bottom-right of a page as decorative
 * background. Sits behind page content (zIndex -1) and is non-interactive, so
 * it only fills otherwise-empty space. The host page must be `position:
 * relative` for the absolute anchoring to work.
 */
import type { SvgIconComponent } from '@mui/icons-material';

export default function PageWatermark({
  icon: Icon,
}: {
  icon: SvgIconComponent;
}) {
  return (
    <Icon
      aria-hidden
      sx={{
        position: 'absolute',
        right: { xs: 4, md: 28 },
        bottom: { xs: 4, md: 24 },
        fontSize: { xs: 200, md: 380 },
        color: 'var(--uk-ac)',
        opacity: 0.08,
        pointerEvents: 'none',
        userSelect: 'none',
        zIndex: -1,
      }}
    />
  );
}
