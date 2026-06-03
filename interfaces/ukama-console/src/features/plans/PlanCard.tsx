/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import Meter from '@/components/Meter';

/**
 * Plan card — color top-bar, price, revenue share (screens-manage.jsx).
 * Shared by Business manage/data-plans (actions) and the agent lens
 * (readOnly) per the §2 invariant.
 */
import { useState } from 'react';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ContentCopyRounded from '@mui/icons-material/ContentCopyRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import EditRounded from '@mui/icons-material/EditRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import { useToast } from '@/components/ToastProvider';
import type { Plan } from '@/data';

export default function PlanCard({
  plan,
  mrr,
  readOnly,
  onEdit,
}: {
  plan: Plan;
  mrr?: number;
  readOnly?: boolean;
  onEdit?: () => void;
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const toast = useToast();
  const rev = plan.subs * plan.price;
  const share = mrr ? Math.round((rev / mrr) * 100) : 0;

  return (
    <div className="card" style={{ padding: 0, overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
      <div style={{ height: 4, background: plan.color }} />
      <div className="card-pad" style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <div style={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}>
            {plan.name}
          </div>
          {!readOnly && (
            <>
              <IconButton
                size="small"
                aria-label="Plan actions"
                sx={{ color: 'var(--uk-ink-3)' }}
                onClick={(e) => setAnchor(e.currentTarget)}
              >
                <MoreVertRounded sx={{ fontSize: 19 }} />
              </IconButton>
              <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
                <MenuItem
                  sx={{ fontSize: 13.5, gap: 1.25 }}
                  onClick={() => {
                    setAnchor(null);
                    onEdit?.();
                  }}
                >
                  <EditRounded sx={{ fontSize: 18 }} /> Edit plan
                </MenuItem>
                <MenuItem
                  sx={{ fontSize: 13.5, gap: 1.25 }}
                  onClick={() => {
                    setAnchor(null);
                    toast(`Duplicated ${plan.name}`);
                  }}
                >
                  <ContentCopyRounded sx={{ fontSize: 18 }} /> Duplicate
                </MenuItem>
                <MenuItem
                  sx={{ fontSize: 13.5, gap: 1.25, color: 'var(--uk-error)' }}
                  onClick={() => {
                    setAnchor(null);
                    toast(`Archived ${plan.name}`, {
                      action: { label: 'Undo', fn: () => toast(`${plan.name} restored`) },
                    });
                  }}
                >
                  <DeleteOutlineRounded sx={{ fontSize: 18 }} /> Archive
                </MenuItem>
              </Menu>
            </>
          )}
        </div>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: 4, margin: '8px 0 2px' }}>
          <span
            className="tnum"
            style={{ fontFamily: 'var(--font-display)', fontSize: 30, fontWeight: 500 }}
          >
            ${plan.price}
          </span>
          <span style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>
            / {plan.days === 1 ? 'day' : 'month'}
          </span>
        </div>
        <div style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>
          {plan.data} data · {plan.days === 1 ? '1 day' : plan.days + ' days'} validity
        </div>
        <hr className="divider" style={{ margin: '14px 0' }} />
        <div style={{ display: 'grid', gap: 9, marginTop: 'auto' }}>
          <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between' }}>
            <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>Customers</span>
            <span
              className="tnum"
              style={{ fontFamily: 'var(--font-display)', fontSize: 18, fontWeight: 500 }}
            >
              {plan.subs}
            </span>
          </div>
          {!readOnly && (
            <>
              <div style={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between' }}>
                <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>Revenue / mo</span>
                <span className="tnum" style={{ fontSize: 14, fontWeight: 600 }}>
                  ${rev.toLocaleString()}
                </span>
              </div>
              <Meter value={share} color={plan.color} sx={{ mt: '2px' }} />
              <div style={{ fontSize: 11, color: 'var(--uk-ink-3)' }}>
                {share}% of plan revenue
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
