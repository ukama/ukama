/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Data plans — PlanCard grid + create (screens-manage.jsx PlansScreen). */
import { useState } from 'react';
import Button from '@mui/material/Button';
import AddRounded from '@mui/icons-material/AddRounded';
import AddCircleOutlineRounded from '@mui/icons-material/AddCircleOutlineRounded';
import PageHeader from '@/components/PageHeader';
import { PLANS } from '@/data';
import CreatePlanDialog from '@/features/plans/CreatePlanDialog';
import PlanCard from '@/features/plans/PlanCard';

export default function PlansScreen() {
  const [showCreate, setShowCreate] = useState(false);
  const mrr = PLANS.reduce((s, p) => s + p.subs * p.price, 0);
  const create = () => setShowCreate(true);

  return (
    <div className="page">
      <PageHeader
        crumb={['Manage', 'Data plans']}
        title="Data plans"
        count={PLANS.length}
        sub="Plans you can assign to customers."
        actions={
          <Button variant="contained" startIcon={<AddRounded />} onClick={create}>
            Create plan
          </Button>
        }
      />
      <div
        className="tile-grid"
        style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))' }}
      >
        {PLANS.map((p) => (
          <PlanCard key={p.id} plan={p} mrr={mrr} onEdit={create} />
        ))}
        <button
          type="button"
          className="card"
          onClick={create}
          style={{
            border: '1.5px dashed var(--uk-line)',
            background: 'transparent',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            gap: 8,
            color: 'var(--uk-ink-3)',
            cursor: 'pointer',
            minHeight: 200,
            fontSize: 14,
            fontWeight: 600,
            fontFamily: 'inherit',
            boxShadow: 'none',
          }}
        >
          <AddCircleOutlineRounded sx={{ fontSize: 30 }} />
          Create a plan
        </button>
      </div>
      {showCreate && <CreatePlanDialog onClose={() => setShowCreate(false)} />}
    </div>
  );
}
