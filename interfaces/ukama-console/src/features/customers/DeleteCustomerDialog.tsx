/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Type-to-confirm delete — friction scaled to risk (findings.jsx
 * DeleteDialog, design finding #11).
 */
import { useState } from 'react';
import Button from '@mui/material/Button';
import DeleteRounded from '@mui/icons-material/DeleteRounded';
import AppModal from '@/components/AppModal';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';

export default function DeleteCustomerDialog({
  sub,
  onClose,
}: {
  sub: Subscriber;
  onClose: () => void;
}) {
  const [val, setVal] = useState('');
  const toast = useToast();
  const ok = val.trim().toLowerCase() === sub.name.toLowerCase();

  return (
    <AppModal
      title="Delete customer"
      width={460}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Cancel
          </Button>
          <Button
            variant="contained"
            color="error"
            startIcon={<DeleteRounded />}
            disabled={!ok}
            onClick={() => {
              onClose();
              toast(`${sub.name} deleted`, {
                action: { label: 'Undo', fn: () => toast(`${sub.name} restored`) },
              });
            }}
          >
            Delete customer
          </Button>
        </>
      }
    >
      <div style={{ fontSize: 14, color: 'var(--uk-ink-2)', lineHeight: 1.55 }}>
        This permanently removes <b style={{ color: 'var(--uk-ink)' }}>{sub.name}</b> and
        releases their SIM ({sub.iccid}) back to the pool. This can’t be undone.
      </div>
      <div style={{ marginTop: 16 }}>
        <label className="flabel">
          Type <b style={{ color: 'var(--uk-ink)' }}>{sub.name}</b> to confirm
        </label>
        <div className="field">
          <input
            value={val}
            onChange={(e) => setVal(e.target.value)}
            placeholder={sub.name}
            autoFocus
          />
        </div>
      </div>
    </AppModal>
  );
}
