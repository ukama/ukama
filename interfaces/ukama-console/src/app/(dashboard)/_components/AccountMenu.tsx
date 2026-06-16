/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Avatar account menu: shows the signed-in user and a Logout action. */
import { useState } from 'react';
import Divider from '@mui/material/Divider';
import ListItemIcon from '@mui/material/ListItemIcon';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import Typography from '@mui/material/Typography';
import LogoutRounded from '@mui/icons-material/LogoutRounded';

import { useAuth } from '@/lib/auth/context';
import { logout } from '@/lib/auth/client';
import { initials } from '@/lib/format';
import { roleLabel } from '@/lib/roles';

export default function AccountMenu() {
  const user = useAuth();
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const [loggingOut, setLoggingOut] = useState(false);

  return (
    <>
      <button
        type="button"
        className="avatar"
        title={user?.name ?? 'Account'}
        aria-haspopup="menu"
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        {initials(user?.name)}
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
        slotProps={{ paper: { sx: { width: 250, mt: 0.5 } } }}
      >
        <div style={{ padding: '8px 16px 6px' }}>
          <Typography sx={{ fontSize: 13.5, fontWeight: 600 }} noWrap>
            {user?.name ?? 'Account'}
          </Typography>
          <Typography sx={{ fontSize: 12, color: 'text.disabled' }} noWrap>
            {user?.email ?? ''}
            {user?.role ? ` · ${roleLabel(user.role)}` : ''}
          </Typography>
        </div>
        <Divider />
        <MenuItem
          sx={{ fontSize: 13.5, gap: 1.25, color: 'var(--uk-error)' }}
          disabled={loggingOut}
          onClick={() => {
            setLoggingOut(true);
            void logout();
          }}
        >
          <ListItemIcon sx={{ color: 'var(--uk-error)' }}>
            <LogoutRounded sx={{ fontSize: 18 }} />
          </ListItemIcon>
          {loggingOut ? 'Signing out…' : 'Log out'}
        </MenuItem>
      </Menu>
    </>
  );
}
