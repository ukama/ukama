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
import Menu from '@mui/material/Menu';
import Typography from '@mui/material/Typography';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import NotificationsRounded from '@mui/icons-material/NotificationsRounded';
import { ALERTS } from '@/data';
import type { Alert } from '@/data';
import { Ic } from './icons';

const SEV_ICON: Record<Alert['sev'], string> = {
  critical: 'error',
  warning: 'warning',
  info: 'info',
};
const SEV_COLOR: Record<Alert['sev'], string> = {
  critical: 'var(--uk-error)',
  warning: 'var(--uk-orange)',
  info: 'var(--uk-ac)',
};

export default function NotificationsMenu() {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const unread = ALERTS.filter((a) => a.sev !== 'info').length;

  return (
    <>
      <button
        type="button"
        className="topbar-icon"
        title="Notifications"
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        <NotificationsRounded sx={{ fontSize: 22 }} />
        {unread > 0 && <span className="badge-num">{unread}</span>}
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 344, mt: 0.5 } } }}
      >
        <Box
          sx={{
            px: 2,
            py: 1,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Typography sx={{ fontWeight: 600, fontSize: 13.5 }}>
            Notifications
          </Typography>
          {unread > 0 && (
            <Typography sx={{ fontSize: 12, color: 'primary.main', fontWeight: 600 }}>
              {unread} unread
            </Typography>
          )}
        </Box>
        {ALERTS.map((a) => (
          <Box
            key={a.id}
            component="button"
            type="button"
            onClick={() => setAnchor(null)} /* TODO Phase 7: open AlertDialog */
            sx={{
              display: 'flex',
              gap: 1.5,
              width: '100%',
              textAlign: 'left',
              alignItems: 'flex-start',
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              px: 2,
              py: 1.25,
              fontFamily: 'inherit',
              '&:hover': { bgcolor: 'action.hover' },
            }}
          >
            <Ic
              name={SEV_ICON[a.sev]}
              sx={{ fontSize: 20, mt: '1px', color: SEV_COLOR[a.sev], flex: 'none' }}
            />
            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Typography sx={{ fontSize: 13, fontWeight: 600 }}>
                {a.title}
              </Typography>
              <Typography
                sx={{ fontSize: 12, color: 'text.secondary', mt: '1px' }}
              >
                {a.detail}
              </Typography>
              <Typography
                sx={{ fontSize: 11.5, color: 'text.disabled', mt: '3px' }}
              >
                {a.site ? `${a.site} · ` : ''}
                {a.age} ago
              </Typography>
            </Box>
            <ChevronRightRounded
              sx={{ fontSize: 18, color: 'text.disabled', alignSelf: 'center', flex: 'none' }}
            />
          </Box>
        ))}
      </Menu>
    </>
  );
}
