/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Network switcher (top bar) — live getNetworks data. The selected id lives
 * in useUiPrefs; if it doesn't resolve to a real network (deleted, or the
 * old seed default) it self-heals to the default/first network. "Add
 * network" opens an inline dialog that creates + selects the new network.
 */
import { useEffect, useState } from 'react';
import Box from '@mui/material/Box';
import ListItemIcon from '@mui/material/ListItemIcon';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';
import AddRounded from '@mui/icons-material/AddRounded';
import CheckRounded from '@mui/icons-material/CheckRounded';
import UnfoldMoreRounded from '@mui/icons-material/UnfoldMoreRounded';

import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { useUiPrefs } from '@/lib/store';
import AddNetworkDialog from './AddNetworkDialog';

export default function NetSwitch() {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const [showAdd, setShowAdd] = useState(false);
  const { networkId, setNetworkId } = useUiPrefs();

  const { data, loading } = useGetNetworksQuery();
  const networks = data?.getNetworks.networks ?? [];
  const selected = networks.find((n) => n.id === networkId);
  const fallback = networks.find((n) => n.isDefault) ?? networks[0];

  // Self-heal a stale selection (seed default / deleted network).
  useEffect(() => {
    if (!loading && !selected && fallback) setNetworkId(fallback.id);
  }, [loading, selected, fallback, setNetworkId]);

  const current = selected ?? fallback;

  return (
    <>
      <button
        type="button"
        className="netswitch"
        onClick={(e) => setAnchor(e.currentTarget)}
        aria-haspopup="menu"
      >
        <span className="dot" />
        <span className="nm">
          {current ? current.name : loading ? '…' : 'No network'}
        </span>
        <UnfoldMoreRounded
          sx={{ fontSize: 18, color: 'rgba(255,255,255,.55)', flexShrink: 0 }}
        />
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 260, mt: 0.5 } } }}
      >
        <Typography
          sx={{
            px: 1.5,
            pt: 1,
            pb: 0.5,
            fontSize: 11,
            fontWeight: 600,
            letterSpacing: '.06em',
            textTransform: 'uppercase',
            color: 'text.disabled',
          }}
        >
          Networks
        </Typography>
        {networks.length === 0 && (
          <Typography
            sx={{ px: 1.5, py: 1, fontSize: 13, color: 'text.disabled' }}
          >
            {loading ? 'Loading networks…' : 'No networks yet'}
          </Typography>
        )}
        {networks.map((n) => (
          <MenuItem
            key={n.id}
            selected={n.id === current?.id}
            onClick={() => {
              setNetworkId(n.id);
              setAnchor(null);
            }}
            sx={{ alignItems: 'flex-start', py: 1 }}
          >
            <Box
              sx={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                mt: '6px',
                mr: 1.25,
                flex: 'none',
                bgcolor: n.isDeactivated
                  ? 'var(--uk-error)'
                  : 'var(--uk-success-bright)',
              }}
            />
            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Typography
                sx={{
                  fontWeight: 600,
                  fontSize: 13.5,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}
                title={n.name}
              >
                {n.name}
              </Typography>
              <Typography sx={{ fontSize: 12, color: 'text.disabled' }}>
                {n.isDefault ? 'Default network' : ' '}
              </Typography>
            </Box>
            {n.id === current?.id && (
              <CheckRounded
                sx={{ fontSize: 18, color: 'primary.main', mt: '4px' }}
              />
            )}
          </MenuItem>
        ))}
        <MenuItem
          onClick={() => {
            setAnchor(null);
            setShowAdd(true);
          }}
          sx={{ mt: 0.5 }}
        >
          <ListItemIcon>
            <AddRounded fontSize="small" />
          </ListItemIcon>
          Add network
        </MenuItem>
      </Menu>
      {showAdd && <AddNetworkDialog onClose={() => setShowAdd(false)} />}
    </>
  );
}
