/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useCallback, useRef, useState } from 'react';

/**
 * Tracks an element's rendered size so the map SVG viewBox can match the
 * container 1:1 (keeps pins/labels at true pixel sizes regardless of the
 * card's aspect ratio).
 */
export function useElementSize(): [
  (el: HTMLElement | null) => void,
  { width: number; height: number },
] {
  const [size, setSize] = useState({ width: 0, height: 0 });
  const obs = useRef<ResizeObserver | null>(null);

  const ref = useCallback((el: HTMLElement | null) => {
    obs.current?.disconnect();
    obs.current = null;
    if (!el) return;
    // measure immediately — ResizeObserver callbacks are throttled in
    // background tabs, so first paint must not depend on them
    const r = el.getBoundingClientRect();
    if (r.width > 0) setSize({ width: r.width, height: r.height });
    obs.current = new ResizeObserver((entries) => {
      const r = entries[0]?.contentRect;
      if (r) {
        setSize((prev) =>
          Math.abs(prev.width - r.width) < 1 && Math.abs(prev.height - r.height) < 1
            ? prev
            : { width: r.width, height: r.height },
        );
      }
    });
    obs.current.observe(el);
  }, []);

  return [ref, size];
}
