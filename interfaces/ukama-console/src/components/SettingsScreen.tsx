/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import TextField from '@mui/material/TextField';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Switch from '@mui/material/Switch';

/**
 * Settings — account / organization / notifications / appearance / billing
 * (screens-manage.jsx SettingsScreen + the v1 appearance controls §7.2 F).
 */
import { useState } from 'react';
import Button from '@mui/material/Button';
import PhotoCameraRounded from '@mui/icons-material/PhotoCameraRounded';
import BillingScreen from '@/app/(dashboard)/business/manage/_components/BillingScreen';
import AppearanceSettings from './AppearanceSettings';
import PageHeader from './PageHeader';
import { useToast } from './ToastProvider';

function Labeled({ label, value }: { label: string; value: string }) {
  return (
    <div className="card card-pad">
      <label className="flabel">{label}</label>
      <TextField fullWidth defaultValue={value} />
    </div>
  );
}

function Toggle({ on: initial }: { on: boolean }) {
  const [on, setOn] = useState(initial);
  return <Switch checked={on} onChange={() => setOn((o) => !o)} />;
}

const TABS = [
  ['account', 'My account'],
  ['org', 'Organization'],
  ['notifs', 'Notifications'],
  ['appearance', 'Appearance'],
  ['billing', 'Billing'],
] as const;

const NOTIF_PREFS: [string, string, boolean][] = [
  ['Site offline alerts', 'Push + email', true],
  ['Low battery warnings', 'Email', true],
  ['Backhaul latency', 'Push', true],
  ['Billing reminders', 'Email', true],
  ['Weekly network digest', 'Email', false],
];

export default function SettingsScreen() {
  const [tab, setTab] = useState<string>('account');
  const toast = useToast();

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
              JM
            </span>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}>
                Joseph Mulenga
              </div>
              <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
                joseph@kwacha.co · Owner
              </div>
            </div>
            <Button
              variant="outlined"
              startIcon={<PhotoCameraRounded />}
              onClick={() => toast('Photo updated')}
            >
              Change photo
            </Button>
          </div>
          <Labeled label="Full name" value="Joseph Mulenga" />
          <Labeled label="Email" value="joseph@kwacha.co" />
          <Labeled label="Phone" value="+260 97 000 1100" />
          <Labeled label="Language" value="English (UK)" />
          <div style={{ gridColumn: '1 / -1', display: 'flex', gap: 10 }}>
            <Button variant="contained" onClick={() => toast('Changes saved')}>
              Save changes
            </Button>
            <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }}>
              Cancel
            </Button>
          </div>
        </div>
      )}

      {tab === 'org' && (
        <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr', maxWidth: 760 }}>
          <Labeled label="Organization name" value="Kwacha Mobile" />
          <Labeled label="Country" value="Zambia" />
          <Labeled label="Currency" value="USD ($)" />
          <Labeled label="Timezone" value="CAT (UTC+2)" />
        </div>
      )}

      {tab === 'notifs' && (
        <div className="card card-pad" style={{ maxWidth: 620 }}>
          {NOTIF_PREFS.map(([l, m, on], i) => (
            <div
              key={l}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '14px 0',
                borderBottom: i < NOTIF_PREFS.length - 1 ? '1px solid var(--uk-line-soft)' : 'none',
              }}
            >
              <div>
                <div style={{ fontSize: 14, fontWeight: 600 }}>{l}</div>
                <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>{m}</div>
              </div>
              <Toggle on={on} />
            </div>
          ))}
        </div>
      )}

      {tab === 'appearance' && (
        <div style={{ maxWidth: 720 }}>
          <AppearanceSettings />
        </div>
      )}

      {tab === 'billing' && <BillingScreen embed />}
    </div>
  );
}
