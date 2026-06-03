/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Customer detail drawer (findings.jsx SubscriberDrawer) — the same record
 * for every lens; actions hidden when readOnly (BUILD-PLAN §2 invariant).
 */
import Button from '@mui/material/Button';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import SwapHorizRounded from '@mui/icons-material/SwapHorizRounded';
import AppDrawer, { DetailRow } from '@/components/AppDrawer';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';

export default function SubscriberDrawer({
  sub,
  onClose,
  readOnly,
}: {
  sub: Subscriber;
  onClose: () => void;
  readOnly?: boolean;
}) {
  const toast = useToast();
  const pct = sub.cap ? Math.min(100, (sub.usage / sub.cap) * 100) : 50;
  const initials = sub.name
    .split(' ')
    .map((x) => x[0])
    .join('');

  return (
    <AppDrawer onClose={onClose} width={430}>
      <div style={{ padding: '20px 24px 16px', borderBottom: '1px solid var(--uk-line)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <div style={{ display: 'flex', gap: 13, alignItems: 'center' }}>
            <span className="av-sm" style={{ width: 46, height: 46, fontSize: 16 }}>
              {initials}
            </span>
            <div>
              <div style={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}>
                {sub.name}
              </div>
              <div className="tnum" style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
                {sub.phone}
              </div>
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close"
            style={{
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              color: 'var(--uk-ink-3)',
              fontSize: 20,
              lineHeight: 1,
              padding: 6,
            }}
          >
            ✕
          </button>
        </div>
      </div>

      <div style={{ flex: 1, overflow: 'auto', padding: '18px 24px' }}>
        <div className="card card-pad" style={{ marginBottom: 14 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
            <span style={{ fontSize: 13, fontWeight: 600 }}>{sub.plan}</span>
            <StatusBadge status={sub.sim === 'suspended' ? 'pending' : sub.sim}>
              {sub.sim === 'suspended' ? 'Suspended' : undefined}
            </StatusBadge>
          </div>
          {sub.cap ? (
            <>
              <div className="meter">
                <span
                  style={{
                    width: pct + '%',
                    background: pct > 90 ? 'var(--uk-orange)' : 'var(--uk-ac)',
                  }}
                />
              </div>
              <div className="tnum" style={{ fontSize: 12.5, color: 'var(--uk-ink-2)', marginTop: 7 }}>
                {sub.usage} of {sub.cap} GB used this cycle
              </div>
            </>
          ) : (
            <div className="tnum" style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}>
              {sub.usage} GB used · unlimited
            </div>
          )}
        </div>

        <DetailRow k="Site" v={sub.site} />
        <DetailRow k="ICCID" v={sub.iccid} />
        <DetailRow k="SIM status" v={<span style={{ textTransform: 'capitalize' }}>{sub.sim}</span>} />
        <DetailRow k="Phone" v={sub.phone} />
      </div>

      {!readOnly && (
        <div
          style={{
            padding: '14px 24px',
            borderTop: '1px solid var(--uk-line)',
            display: 'flex',
            gap: 10,
          }}
        >
          <Button
            variant="contained"
            startIcon={<AddCardRounded />}
            sx={{ flex: 1 }}
            onClick={() => toast(`Top up for ${sub.name} — flow lands with the form dialogs`)}
          >
            Top up
          </Button>
          <Button
            variant="outlined"
            startIcon={<SwapHorizRounded />}
            sx={{ flex: 1 }}
            onClick={() => toast(`Change plan for ${sub.name} — flow lands with the form dialogs`)}
          >
            Change plan
          </Button>
        </div>
      )}
    </AppDrawer>
  );
}
