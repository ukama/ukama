/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import DarkModeRounded from '@mui/icons-material/DarkModeRounded';
import LightModeRounded from '@mui/icons-material/LightModeRounded';
import MenuRounded from '@mui/icons-material/MenuRounded';
import { useColorScheme } from '@mui/material/styles';
import { useSyncExternalStore } from 'react';

import { useAuth } from '@/lib/auth/context';
import AccountMenu from './AccountMenu';
import LensSegment from './LensSegment';
import NetSwitch from './NetSwitch';
import NotificationsMenu from './NotificationsMenu';

const emptySubscribe = () => () => {};
const useMounted = () =>
  useSyncExternalStore(
    emptySubscribe,
    () => true,
    () => false,
  );

function ThemeToggle() {
  const { mode, setMode, systemMode } = useColorScheme();
  const mounted = useMounted();
  if (!mounted) return <span className="topbar-icon" />;
  const dark = (mode === 'system' ? systemMode : mode) === 'dark';
  return (
    <button
      type="button"
      className="topbar-icon"
      title={dark ? 'Switch to light mode' : 'Switch to dark mode'}
      onClick={() => setMode(dark ? 'light' : 'dark')}
    >
      {dark ? (
        <LightModeRounded sx={{ fontSize: 21 }} />
      ) : (
        <DarkModeRounded sx={{ fontSize: 21 }} />
      )}
    </button>
  );
}

/** Authenticated user's org name — plain label, bold, theme accent color. */
function OrgLabel() {
  const user = useAuth();
  if (!user?.orgName) return null;
  return (
    <span
      title={`Organization: ${user.orgName}`}
      style={{
        fontFamily: 'var(--font-display)',
        fontSize: 14,
        fontWeight: 700,
        color: 'var(--uk-ac)',
        whiteSpace: 'nowrap',
        marginRight: 4,
      }}
    >
      {user.orgName}
    </span>
  );
}

export default function TopBar({ onMenu }: { onMenu: () => void }) {
  return (
    <header className="topbar">
      <button
        type="button"
        className="topbar-icon menu-btn"
        aria-label="Open navigation"
        onClick={onMenu}
      >
        <MenuRounded sx={{ fontSize: 22 }} />
      </button>
      <NetSwitch />
      <LensSegment />
      <div className="spacer" />
      <OrgLabel />
      <ThemeToggle />
      <NotificationsMenu />
      <AccountMenu />
    </header>
  );
}
