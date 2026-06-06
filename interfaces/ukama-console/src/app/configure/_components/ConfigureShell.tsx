/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * /configure chrome: brand, STEP x/y, right-rail progressive hierarchy, and
 * resume-point persistence — every step change records the current URL in
 * useUiPrefs.lastConfigureUrl ("Continue setup" lands back here). The
 * complete step clears it.
 */
'use client';

import { useEffect } from 'react';
import { usePathname } from 'next/navigation';
import PublicRounded from '@mui/icons-material/PublicRounded';
import LanRounded from '@mui/icons-material/LanRounded';
import CellTowerRounded from '@mui/icons-material/CellTowerRounded';
import SimCardRounded from '@mui/icons-material/SimCardRounded';

import UMark from '@/components/UMark';
import { useAuth } from '@/lib/auth/context';
import { useUiPrefs } from '@/lib/store';
import { useCurrentConfigureUrl } from './state';

/** Ordered steps with a STEP x/y label (overview/complete have none). */
const STEPS = [
  '/configure/network',
  '/configure/install',
  '/configure/site',
  '/configure/site/settings',
  '/configure/sims',
] as const;

/** Index of the active step by exact path match (-1 if not a numbered step). */
const stepIndexFor = (pathname: string): number =>
  STEPS.reduce((best, step, i) => (pathname === step ? i : best), -1);

function TreeNode({
  icon,
  label,
  active,
  last,
}: {
  icon: React.ReactNode;
  label: string;
  active: boolean;
  last?: boolean;
}) {
  return (
    <div className="cfg-tree-node" data-active={active}>
      <div className="cfg-tree-icon">{icon}</div>
      <div className="cfg-tree-label">{label}</div>
      {!last && <div className="cfg-tree-connector" />}
    </div>
  );
}

export default function ConfigureShell({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const currentUrl = useCurrentConfigureUrl();
  const setLastConfigureUrl = useUiPrefs((s) => s.setLastConfigureUrl);
  const user = useAuth();
  const orgName = user?.orgName ?? 'Organization';

  const stepIndex = stepIndexFor(pathname);
  const isComplete = pathname.startsWith('/configure/complete');

  // Persist the resume point on every step change; clear it on completion.
  useEffect(() => {
    setLastConfigureUrl(isComplete ? null : currentUrl);
  }, [currentUrl, isComplete, setLastConfigureUrl]);

  return (
    <main className="cfg-root">
      <section className="cfg-main">
        <div className="cfg-brand">
          <span style={{ display: 'inline-flex', width: 22 }}>
            <UMark />
          </span>
          <span className="cfg-brand-name">ukama</span>
        </div>

        <div
          className="cfg-step-label"
          style={{ visibility: stepIndex >= 0 ? 'visible' : 'hidden' }}
        >
          Step {stepIndex + 1}/{STEPS.length}
        </div>

        <div className="cfg-body">{children}</div>
      </section>

      <aside className="cfg-aside" aria-hidden="true">
        <div>
          <TreeNode
            icon={<PublicRounded sx={{ fontSize: 24 }} />}
            label={orgName}
            active
          />
          <TreeNode
            icon={<LanRounded sx={{ fontSize: 24 }} />}
            label="Network"
            active={stepIndex >= 0 || isComplete}
          />
          <TreeNode
            icon={<CellTowerRounded sx={{ fontSize: 24 }} />}
            label="Site"
            active={stepIndex >= 1 || isComplete}
          />
          <TreeNode
            icon={<SimCardRounded sx={{ fontSize: 24 }} />}
            label="SIMs"
            active={stepIndex >= 4 || isComplete}
            last
          />
        </div>
      </aside>
    </main>
  );
}
