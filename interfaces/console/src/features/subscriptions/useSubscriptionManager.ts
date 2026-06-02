/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { useCallback, useEffect, useRef } from 'react';

export const useSubscriptionManager = () => {
  const activeTopics = useRef<Map<string, () => void>>(new Map());

  const subscribe = useCallback((topic: string, cleanup: () => void) => {
    if (activeTopics.current.has(topic)) return;
    activeTopics.current.set(topic, cleanup);
  }, []);

  const unsubscribe = useCallback((topic: string) => {
    const cleanup = activeTopics.current.get(topic);
    cleanup?.();
    activeTopics.current.delete(topic);
  }, []);

  // Auto-cleanup all on unmount
  useEffect(() => {
    const topics = activeTopics.current;
    return () => {
      topics.forEach((cleanup) => cleanup());
      topics.clear();
    };
  }, []);

  return { subscribe, unsubscribe };
};
