/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * First-visit welcome screen: confirms org, country, and member role before
 * entering the console. "Continue" acknowledges the welcome via
 * /api/auth/welcome (BFF records it, token cookie is cleared) and performs a
 * full navigation so the proxy re-mints a token with isShowWelcome=false.
 * Layout/responsive rules live in ../welcome.css.
 */
'use client';

import { useState } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import PublicRounded from '@mui/icons-material/PublicRounded';
import ApartmentRounded from '@mui/icons-material/ApartmentRounded';
import PersonRounded from '@mui/icons-material/PersonRounded';

import UMark from '@/components/UMark';
import { resolveResumeUrl, useActivation } from '@/lib/activation';
import type { AuthUser } from '@/lib/auth/types';
import { useUiPrefs } from '@/lib/store';

/** Token roles (ROLE_*) → display labels used across the console. */
const ROLE_LABELS: Record<string, string> = {
  ROLE_OWNER: 'Owner',
  ROLE_ADMIN: 'Administrator',
  ROLE_NETWORK_OWNER: 'Network owner',
  ROLE_VENDOR: 'Vendor',
  ROLE_USER: 'Member',
};

function roleLabel(role: string): string {
  if (ROLE_LABELS[role]) return ROLE_LABELS[role];
  const cleaned = role.replace(/^ROLE_/, '').replace(/_/g, ' ').toLowerCase();
  return cleaned ? cleaned[0]?.toUpperCase() + cleaned.slice(1) : role;
}

/**
 * The token carries an ISO country code while a fresh /get-user response
 * carries the full name — resolve codes via Intl and pass names through.
 */
function countryLabel(country: string): string {
  if (/^[A-Z]{2}$/.test(country)) {
    try {
      return (
        new Intl.DisplayNames(['en'], { type: 'region' }).of(country) ??
        country
      );
    } catch {
      return country;
    }
  }
  return country;
}

function Field({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="welcome-field-label">{label}</div>
      <div className="welcome-field-value">{value}</div>
    </div>
  );
}

function TreeNode({
  icon,
  label,
  last,
}: {
  icon: React.ReactNode;
  label: string;
  last?: boolean;
}) {
  return (
    <div className="welcome-tree-node">
      <div className="welcome-tree-icon">{icon}</div>
      <div className="welcome-tree-label">{label}</div>
      {!last && <div className="welcome-tree-connector" />}
    </div>
  );
}

export default function WelcomeScreen({ user }: { user: AuthUser }) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  // Pre-fetches activation state while the user reads the page, so Continue
  // can land first-time admins directly in /configure (unknown state → '/').
  const { status, needsSetup } = useActivation();
  const lastConfigureUrl = useUiPrefs((s) => s.lastConfigureUrl);

  const onContinue = async () => {
    setBusy(true);
    setError(null);
    try {
      const res = await fetch('/api/auth/welcome', { method: 'POST' });
      if (!res.ok) throw new Error(`acknowledge failed (${res.status})`);
      // Full navigation (not router.push) so proxy.ts re-mints the session
      // token — the fresh token carries isShowWelcome=false.
      window.location.assign(
        needsSetup ? resolveResumeUrl(status, lastConfigureUrl) : '/',
      );
    } catch {
      setError("Couldn't save your confirmation. Please try again.");
      setBusy(false);
    }
  };

  const memberLine = `${user.name} | ${roleLabel(user.role)}`;

  return (
    <main className="welcome-root">
      {/* Left: welcome copy + details */}
      <section className="welcome-main">
        <div className="welcome-brand">
          <span style={{ display: 'inline-flex', width: 22 }}>
            <UMark />
          </span>
          <span className="welcome-brand-name">ukama</span>
        </div>

        <div className="welcome-body">
          <h1 className="welcome-title">Welcome to Ukama!</h1>
          <p className="welcome-copy">
            Please check to make sure the following details are correct before
            continuing to the Console, where you can manage and monitor your
            network.
          </p>

          <div className="welcome-fields">
            <Field
              label="Network operating country"
              value={countryLabel(user.country)}
            />
            <Field label="Organization name" value={user.orgName} />
            <Field label="Role" value={memberLine} />
          </div>
        </div>

        <div className="welcome-actions">
          {error && <p className="welcome-error">{error}</p>}
          <Button
            variant="contained"
            size="large"
            onClick={() => void onContinue()}
            disabled={busy}
            startIcon={
              busy ? <CircularProgress size={16} color="inherit" /> : undefined
            }
          >
            {busy ? 'Continuing…' : 'Continue'}
          </Button>
        </div>
      </section>

      {/* Right: org hierarchy summary (hidden on mobile — see welcome.css) */}
      <aside className="welcome-aside">
        <div>
          <TreeNode
            icon={<PublicRounded sx={{ fontSize: 26 }} />}
            label={countryLabel(user.country)}
          />
          <TreeNode
            icon={<ApartmentRounded sx={{ fontSize: 26 }} />}
            label={user.orgName}
          />
          <TreeNode
            icon={<PersonRounded sx={{ fontSize: 26 }} />}
            label={memberLine}
            last
          />
        </div>
      </aside>
    </main>
  );
}
