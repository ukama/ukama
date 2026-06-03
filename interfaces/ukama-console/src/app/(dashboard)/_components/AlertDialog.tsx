/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Alert action sheet from the notifications feed (findings.jsx AlertDialog). */
import Button from '@mui/material/Button';
import BuildRounded from '@mui/icons-material/BuildRounded';
import PlaceRounded from '@mui/icons-material/PlaceRounded';
import AppModal from '@/components/AppModal';
import type { Alert } from '@/data';
import { Ic } from './icons';

export default function AlertDialog({
  alert,
  onClose,
  onAction,
}: {
  alert: Alert;
  onClose: () => void;
  onAction: (alert: Alert) => void;
}) {
  return (
    <AppModal
      title={alert.title}
      width={460}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Dismiss
          </Button>
          <Button variant="contained" startIcon={<BuildRounded />} onClick={() => onAction(alert)}>
            {alert.action}
          </Button>
        </>
      }
    >
      <div className={`alert-row sev-${alert.sev}`}>
        <span className="alert-ic">
          <Ic name={alert.icon} sx={{ fontSize: 19 }} />
        </span>
        <div>
          <div style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.5 }}>
            {alert.detail}
          </div>
          {alert.site && (
            <div
              style={{
                fontSize: 12.5,
                color: 'var(--uk-ink-3)',
                marginTop: 6,
                display: 'flex',
                alignItems: 'center',
                gap: 3,
              }}
            >
              <PlaceRounded sx={{ fontSize: 14 }} /> {alert.site} · {alert.age} ago
            </div>
          )}
        </div>
      </div>
      <div
        className="card card-pad"
        style={{
          marginTop: 16,
          background: 'var(--uk-page)',
          border: 'none',
          boxShadow: 'none',
          fontSize: 13,
          color: 'var(--uk-ink-2)',
        }}
      >
        <b style={{ color: 'var(--uk-ink)' }}>Suggested checks</b>
        <ul style={{ margin: '8px 0 0', paddingLeft: 18, lineHeight: 1.7 }}>
          <li>Confirm power source and battery state</li>
          <li>Ping the node’s backhaul link</li>
          <li>Review last config change</li>
        </ul>
      </div>
    </AppModal>
  );
}
