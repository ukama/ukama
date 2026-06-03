/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useState } from 'react';
import Box from '@mui/material/Box';
import ListItemIcon from '@mui/material/ListItemIcon';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';
import AddRounded from '@mui/icons-material/AddRounded';
import CheckRounded from '@mui/icons-material/CheckRounded';
import UnfoldMoreRounded from '@mui/icons-material/UnfoldMoreRounded';
import { NETWORKS } from '@/data';
import { useUiPrefs } from '@/lib/store';
import OnboardingFlow from './OnboardingFlow';

const STATUS_DOT: Record<string, string> = {
  online: 'var(--uk-success-bright)',
  degraded: 'var(--uk-warning)',
  offline: 'var(--uk-error)',
};

export default function NetSwitch() {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const [showOnboarding, setShowOnboarding] = useState(false);
  const { networkId, setNetworkId } = useUiPrefs();
  const net = NETWORKS.find((n) => n.id === networkId) ?? NETWORKS[0];

  return (
    <>
      <button
        type="button"
        className="netswitch"
        onClick={(e) => setAnchor(e.currentTarget)}
        aria-haspopup="menu"
      >
        <span className="dot" />
        <span className="nm">{net ? net.name : 'Network'}</span>
        <UnfoldMoreRounded
          sx={{ fontSize: 18, color: 'rgba(255,255,255,.55)' }}
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
        {NETWORKS.map((n) => (
          <MenuItem
            key={n.id}
            selected={n.id === networkId}
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
                bgcolor: STATUS_DOT[n.status],
              }}
            />
            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Typography sx={{ fontWeight: 600, fontSize: 13.5 }}>
                {n.name}
              </Typography>
              <Typography sx={{ fontSize: 12, color: 'text.disabled' }}>
                {n.region}
              </Typography>
            </Box>
            {n.id === networkId && (
              <CheckRounded sx={{ fontSize: 18, color: 'primary.main', mt: '4px' }} />
            )}
          </MenuItem>
        ))}
        <MenuItem
          onClick={() => {
            setAnchor(null);
            setShowOnboarding(true);
          }}
          sx={{ mt: 0.5 }}
        >
          <ListItemIcon>
            <AddRounded fontSize="small" />
          </ListItemIcon>
          Add network
        </MenuItem>
      </Menu>
      {showOnboarding && <OnboardingFlow onClose={() => setShowOnboarding(false)} />}
    </>
  );
}
