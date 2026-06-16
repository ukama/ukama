/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * App-level banner shown above the top bar when one or more nodes are online
 * but still in the "unknown" state (reachable but not configured). "Configure
 * node" opens the node's detail page. Dismissible for the session.
 */
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import CloseRounded from '@mui/icons-material/CloseRounded';
import InfoOutlined from '@mui/icons-material/InfoOutlined';

import { useNodePoolQuery } from '@/client/graphql/nodes-list.generated';

const DISMISS_KEY = 'uk-node-config-bar-dismissed';

export default function UnconfiguredNodeBar() {
  const router = useRouter();
  // Pool query (all nodes, including unassigned/available) so the prompt
  // covers brand-new nodes not yet attached to a network.
  const { data } = useNodePoolQuery();
  const [dismissed, setDismissed] = useState(
    () =>
      typeof window !== 'undefined' &&
      sessionStorage.getItem(DISMISS_KEY) === '1',
  );

  // Online + "unknown" state = reachable but not configured yet. Offline
  // nodes can't be configured, so they don't prompt.
  const pending = (data?.nodesView.nodes.nodes ?? []).filter(
    (n) =>
      n.status.connectivity?.toLowerCase() === 'online' &&
      n.status.state?.toLowerCase() === 'unknown',
  );
  if (dismissed || pending.length === 0) return null;

  const message =
    pending.length === 1
      ? 'New node is available and ready to configure.'
      : `${pending.length} new nodes are available and ready to configure.`;

  return (
    <div className="setup-bar" role="status">
      <InfoOutlined sx={{ fontSize: 18 }} />
      <span className="setup-bar-text">{message}</span>
      <Button
        size="small"
        variant="contained"
        disableElevation
        onClick={() => router.push('/configure/install')}
      >
        Configure node
      </Button>
      <IconButton
        size="small"
        aria-label="Dismiss"
        onClick={() => {
          sessionStorage.setItem(DISMISS_KEY, '1');
          setDismissed(true);
        }}
        sx={{ color: 'inherit' }}
      >
        <CloseRounded sx={{ fontSize: 18 }} />
      </IconButton>
    </div>
  );
}
