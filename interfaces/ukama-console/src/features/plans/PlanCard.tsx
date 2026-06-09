/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Plan card — color top-bar, price, validity, and customer count. Shared by
 * Business manage/data-plans (with an Edit action) and the agent lens
 * (readOnly, no actions).
 */
import { useState } from 'react';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import EditRounded from '@mui/icons-material/EditRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import type { Plan } from '@/data';
import { useCurrency } from '@/lib/currency';

export default function PlanCard({
  plan,
  readOnly,
  onEdit,
}: {
  plan: Plan;
  readOnly?: boolean;
  onEdit?: () => void;
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const { symbol } = useCurrency();

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
              </Menu>
            </>
          )}
        </div>
        <div style={{ display: 'flex', alignItems: 'baseline', gap: 4, margin: '8px 0 2px' }}>
          <span
            className="tnum"
            style={{ fontFamily: 'var(--font-display)', fontSize: 30, fontWeight: 500 }}
          >
            {symbol}{plan.price}
          </span>
          <span style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>
            / {plan.days === 1 ? 'day' : 'month'}
          </span>
        </div>
        <div style={{ fontSize: 13.5, color: 'var(--uk-ink-2)' }}>
          {plan.data} data · {plan.days === 1 ? '1 day' : plan.days + ' days'} validity
        </div>
      </div>
    </div>
  );
}
