/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * First-visit load shimmer: a screen skeletons the first time it's opened in
 * a session, then never again — navigating back is instant (table-kit.jsx).
 * Mirrors the §5.1 cache-first behaviour until Apollo lands in Phase 2 (API).
 */
import { useEffect, useState } from 'react';

const seen = new Set<string>();

export function useFirstLoad(key: string, ms = 620): boolean {
  const [loading, setLoading] = useState(() => !seen.has(key));

  useEffect(() => {
    if (seen.has(key)) return;
    const t = setTimeout(() => {
      seen.add(key);
      setLoading(false);
    }, ms);
    return () => clearTimeout(t);
  }, [key, ms]);

  return loading;
}
