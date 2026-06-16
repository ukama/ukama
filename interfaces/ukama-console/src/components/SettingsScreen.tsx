/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Settings — account / organization / billing, bound to the signed-in user. */
import { useState } from 'react';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import { useAuth } from '@/lib/auth/context';
import { useCurrency } from '@/lib/currency';
import { countryLabel, initials } from '@/lib/format';
import { roleLabel } from '@/lib/roles';
import PageHeader from './PageHeader';

function ReadField({ label, value }: { label: string; value: string }) {
  return (
    <div className="card card-pad">
      <label className="flabel">{label}</label>
      <div className="ff-readonly">{value || '—'}</div>
    </div>
  );
}

const TABS = [
  ['account', 'My account'],
  ['org', 'Organization'],
] as const;

export default function SettingsScreen() {
  const [tab, setTab] = useState<string>('account');
  const user = useAuth();
  const { symbol, code } = useCurrency();

  const name = user?.name ?? '—';
  const email = user?.email ?? '—';
  const role = user ? roleLabel(user.role) : '—';

  return (
    <div className="page">
      <PageHeader title="Settings" sub="Manage your account, organization and billing." />
      <Tabs value={tab} onChange={(_, v: string) => setTab(v)}>
        {TABS.map(([k, l]) => (
          <Tab key={k} value={k} label={l} />
        ))}
      </Tabs>

      {tab === 'account' && (
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', maxWidth: 760 }}>
          <div
            className="card card-pad"
            style={{ gridColumn: '1 / -1', display: 'flex', alignItems: 'center', gap: 18 }}
          >
            <span className="av-sm" style={{ width: 60, height: 60, fontSize: 22 }}>
              {initials(name, email)}
            </span>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}>
                {name}
              </div>
              <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
                {email} · {role}
              </div>
            </div>
          </div>
          <ReadField label="Full name" value={name} />
          <ReadField label="Email" value={email} />
          <ReadField label="Role" value={role} />
          <ReadField
            label="Email verified"
            value={user?.isEmailVerified ? 'Yes' : 'No'}
          />
        </div>
      )}

      {tab === 'org' && (
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', maxWidth: 760 }}>
          <ReadField label="Organization name" value={user?.orgName ?? '—'} />
          <ReadField label="Country" value={countryLabel(user?.country ?? '')} />
          <ReadField
            label="Currency"
            value={code ? `${code}${symbol ? ` (${symbol})` : ''}` : '—'}
          />
        </div>
      )}
    </div>
  );
}
