/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import { useMemo, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import Box from '@mui/material/Box';
import Dialog from '@mui/material/Dialog';
import Divider from '@mui/material/Divider';
import InputBase from '@mui/material/InputBase';
import Typography from '@mui/material/Typography';
import SearchRounded from '@mui/icons-material/SearchRounded';
import { NAV_BY_LENS, bottomNav, lensFromPath } from '../_config/nav';
import { Ic } from './icons';

/**
 * Inner content holds the query state and is mounted fresh on every open
 * (Dialog unmounts children when closed), so the search resets naturally.
 */
function PaletteContent({ onClose }: { onClose: () => void }) {
  const [q, setQ] = useState('');
  const pathname = usePathname();
  const router = useRouter();
  const lens = lensFromPath(pathname);

  const pages = useMemo(() => {
    const items = NAV_BY_LENS[lens].flatMap((g) => g.items);
    return [...items, ...bottomNav(lens)];
  }, [lens]);

  const ql = q.toLowerCase();
  const fp = pages.filter((p) => p.label.toLowerCase().includes(ql));

  const go = (href: string) => {
    router.push(href);
    onClose();
  };

  const rowSx = {
    display: 'flex',
    alignItems: 'center',
    gap: 1.5,
    width: '100%',
    border: 'none',
    background: 'transparent',
    cursor: 'pointer',
    textAlign: 'left' as const,
    px: 1.5,
    py: 1,
    borderRadius: 1.5,
    fontFamily: 'inherit',
    '&:hover': { bgcolor: 'action.hover' },
  };

  return (
    <>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.25, px: 2, py: 1.5 }}>
        <SearchRounded sx={{ fontSize: 22, color: 'text.disabled' }} />
        <InputBase
          autoFocus
          fullWidth
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Jump to a page…"
          sx={{ fontSize: 16 }}
        />
        <Typography
          sx={{
            border: '1px solid',
            borderColor: 'divider',
            color: 'text.disabled',
            px: 0.75,
            py: 0.25,
            borderRadius: 1,
            fontSize: 11,
          }}
        >
          esc
        </Typography>
      </Box>
      <Divider />
      <Box sx={{ maxHeight: '48vh', overflow: 'auto', p: 1 }}>
        {fp.length > 0 && (
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
            Pages
          </Typography>
        )}
        {fp.map((p) => (
          <Box key={p.href} component="button" type="button" sx={rowSx} onClick={() => go(p.href)}>
            <Ic name={p.icon} sx={{ fontSize: 20, color: 'text.disabled' }} />
            <Typography sx={{ fontSize: 13.5 }}>{p.label}</Typography>
          </Box>
        ))}
        {fp.length === 0 && (
          <Typography sx={{ p: 3.5, textAlign: 'center', color: 'text.disabled' }}>
            No matches for “{q}”.
          </Typography>
        )}
      </Box>
    </>
  );
}

export default function CommandPalette({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      slotProps={{
        paper: {
          sx: { width: 560, mt: '12vh', alignSelf: 'flex-start', borderRadius: 3 },
        },
      }}
      sx={{ '& .MuiDialog-container': { alignItems: 'flex-start' } }}
    >
      <PaletteContent onClose={onClose} />
    </Dialog>
  );
}
