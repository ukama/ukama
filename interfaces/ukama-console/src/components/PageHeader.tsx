/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Page header — crumb · title (+count) · sub · "Last fetched HH:MM" · actions.
 * Ported from the prototype's PageHead (shell.jsx); shared by every screen.
 */
import { useState } from 'react';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import ScheduleRounded from '@mui/icons-material/ScheduleRounded';

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
  const t = fetchedAt.toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div className="pagehead">
      <div>
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
        <div className="pagetitle">
          {title}
          {count != null && <span className="cnt tnum">{count}</span>}
        </div>
        {sub && <div className="pagesub">{sub}</div>}
      </div>
      <div className="pagehead-right">
        <span className="page-fetched" suppressHydrationWarning>
          <ScheduleRounded sx={{ fontSize: 15 }} />
          Last fetched {t}
        </span>
        {actions && <div className="head-actions">{actions}</div>}
      </div>
    </div>
  );
}
