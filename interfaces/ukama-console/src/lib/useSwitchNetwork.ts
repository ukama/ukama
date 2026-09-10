/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Selects a network on the user's behalf (the top-bar switcher, "Add
 * network"). A site or node detail page belongs to the network it was opened
 * from, so switching away from one lands on the matching list for the new
 * network instead of leaving the old entity on screen.
 *
 * NetSwitch's self-heal of a stale stored id calls `setNetworkId` directly,
 * not this hook: it is not a user switch, and redirecting there would bounce
 * a deep link to a detail page back to its list.
 */
import { useCallback } from 'react';
import { usePathname, useRouter } from 'next/navigation';

import { useUiPrefs } from '@/lib/store';

const DETAIL_TO_LIST: [RegExp, string][] = [
  [/^\/network\/sites\/[^/]+/, '/network/sites'],
  [/^\/network\/nodes\/[^/]+/, '/network/nodes'],
];

export function useSwitchNetwork(): (networkId: string) => void {
  const router = useRouter();
  const pathname = usePathname();
  const currentId = useUiPrefs((s) => s.networkId);
  const setNetworkId = useUiPrefs((s) => s.setNetworkId);

  return useCallback(
    (networkId: string) => {
      if (networkId === currentId) return;

      setNetworkId(networkId);

      const list = DETAIL_TO_LIST.find(([re]) => re.test(pathname))?.[1];
      // replace: Back should not return to the previous network's entity.
      if (list) router.replace(list);
    },
    [currentId, pathname, router, setNetworkId],
  );
}
