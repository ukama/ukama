/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import ChevronLeftRounded from '@mui/icons-material/ChevronLeftRounded';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import { useAuth } from '@/lib/auth/context';
import { publicEnv } from '@/lib/runtime-env';
import { useUiPrefs } from '@/lib/store';
import { NAV_BY_LENS, bottomNav, lensFromPath } from '../_config/nav';
import type { NavItem } from '../_config/nav';
import { Ic } from './icons';

function isActive(pathname: string, item: NavItem): boolean {
  if (item.exact) return pathname === item.href;
  return pathname === item.href || pathname.startsWith(item.href + '/');
}

function NavLink({ item, pathname }: { item: NavItem; pathname: string }) {
  const active = isActive(pathname, item);
  return (
    <Link
      href={item.href}
      className={`navitem${active ? ' active' : ''}`}
      title={item.label}
    >
      <Ic name={item.icon} className="ni-ic" />
      <span className="ni-label">{item.label}</span>
    </Link>
  );
}

/**
 * External sidebar link (e.g. "Backend status"): opens the status app in a
 * new tab with the session's org passed via ?org= so the status app can
 * scope its init lookups. No active-state matching — it never owns a
 * console route.
 */
function ExternalNavLink({ item }: { item: NavItem }) {
  const user = useAuth();
  const href = `${publicEnv().statusAppUrl}/?org=${encodeURIComponent(
    user?.orgName ?? '',
  )}`;
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="navitem"
      title={item.label}
    >
      <Ic name={item.icon} className="ni-ic" />
      <span className="ni-label">{item.label}</span>
    </a>
  );
}

export default function Sidebar() {
  const pathname = usePathname();
  const lens = lensFromPath(pathname);
  const groups = NAV_BY_LENS[lens];
  const { rail, toggleRail } = useUiPrefs();

  return (
    <aside className="sidebar">
      {groups.map((g, gi) => (
        <div className="nav-group" key={gi}>
          {g.group && <div className="nav-label">{g.group}</div>}
          {g.items.map((item) => (
            <NavLink key={item.href} item={item} pathname={pathname} />
          ))}
        </div>
      ))}
      <div className="grow" />
      <hr className="sidebar-divider" />
      {bottomNav(lens).map((item) =>
        item.external ? (
          <ExternalNavLink key={item.href} item={item} />
        ) : (
          <NavLink key={item.href} item={item} pathname={pathname} />
        ),
      )}
      <button type="button" className="railtoggle" onClick={toggleRail}>
        {rail === 'icon' ? (
          <ChevronRightRounded sx={{ fontSize: 20 }} />
        ) : (
          <ChevronLeftRounded sx={{ fontSize: 20 }} />
        )}
        <span className="ni-label">Collapse</span>
      </button>
    </aside>
  );
}
