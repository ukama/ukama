/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Page header — crumb · title (+count) · sub · "Last update: time · date" · actions.
 * Ported from the prototype's PageHead (shell.jsx); shared by every screen.
 */
import { useEffect, useState } from 'react';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import ScheduleRounded from '@mui/icons-material/ScheduleRounded';

/** Compact "time ago" label relative to now. */
function timeAgo(from: Date, now: Date): string {
  const s = Math.max(0, Math.round((now.getTime() - from.getTime()) / 1000));
  if (s < 10) return 'just now';
  if (s < 60) return `${s}s ago`;
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  return `${Math.floor(h / 24)}d ago`;
}

export interface PageHeaderProps {
  crumb?: string[];
  title: string;
  count?: string | number;
  sub?: string;
  actions?: React.ReactNode;
}

export default function PageHeader({
  crumb,
  title,
  count,
  sub,
  actions,
}: PageHeaderProps) {
  const [fetchedAt] = useState(() => new Date());
  const [now, setNow] = useState(fetchedAt);
  // Re-tick the relative label without re-fetching anything.
  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 10_000);
    return () => clearInterval(id);
  }, []);
  const rel = timeAgo(fetchedAt, now);
  const full =
    fetchedAt.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      second: '2-digit',
      hour12: true,
    }) +
    ' \u00b7 ' +
    fetchedAt.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });

  // Flat children so the grid can align rows: crumb\u2194timestamp, title\u2194actions.
  return (
    <div className={`pagehead${crumb ? ' has-crumb' : ''}`}>
      {crumb && (
        <div className="crumb">
          {crumb.map((c, i) => (
            <span key={i} style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
              {i > 0 && <ChevronRightRounded sx={{ fontSize: 15 }} />}
              <span>{c}</span>
            </span>
          ))}
        </div>
      )}
      <span className="page-fetched" title={`Last update: ${full}`} suppressHydrationWarning>
        <ScheduleRounded sx={{ fontSize: 15 }} />
        Updated {rel}
      </span>
      <div className="pagetitle">
        {title}
        {/* Show the count only when meaningful (a non-zero number, or a
            custom string like a range). Zero/empty is noise. */}
        {count != null &&
          !(typeof count === 'number' && count === 0) &&
          count !== '' && <span className="cnt tnum">{count}</span>}
      </div>
      {actions && <div className="head-actions">{actions}</div>}
      {sub && <div className="pagesub">{sub}</div>}
    </div>
  );
}
