/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Agent data plans — read-only browse with search + sort (screens-manage.jsx). */
import { useState } from 'react';
import { EmptyState } from '@/components/EmptyState';
import PageHeader from '@/components/PageHeader';
import SearchField from '@/components/SearchField';
import { PLANS } from '@/data';
import PlanCard from '@/features/plans/PlanCard';

type Sort = 'price-asc' | 'price-desc' | 'data-desc';

export default function AgentPlansScreen() {
  const [q, setQ] = useState('');
  const [sort, setSort] = useState<Sort>('price-asc');
  const vol = (data: string) =>
    /unlim/i.test(data) ? Infinity : parseFloat(data) || 0;

  let list = PLANS.filter((p) => p.name.toLowerCase().includes(q.toLowerCase()));
  list = [...list].sort((a, b) =>
    sort === 'price-asc'
      ? a.price - b.price
      : sort === 'price-desc'
        ? b.price - a.price
        : vol(b.data) - vol(a.data),
  );

  return (
    <div className="page">
      <PageHeader
        title="Data plans"
        count={PLANS.length}
        sub="Browse plans to assign, top up or change for your customers."
      />
      <div style={{ display: 'flex', gap: 10, marginBottom: 18, flexWrap: 'wrap', alignItems: 'center' }}>
        <SearchField value={q} onChange={setQ} placeholder="Search plans" />
        <div className="seg">
          {(
            [
              ['price-asc', 'Price ↑'],
              ['price-desc', 'Price ↓'],
              ['data-desc', 'Data'],
            ] as const
          ).map(([k, l]) => (
            <button
              key={k}
              type="button"
              className={sort === k ? 'on' : ''}
              onClick={() => setSort(k)}
            >
              {l}
            </button>
          ))}
        </div>
        <div style={{ marginLeft: 'auto', fontSize: 13, color: 'var(--uk-ink-3)' }}>
          {list.length} of {PLANS.length}
        </div>
      </div>
      {list.length === 0 ? (
        <div className="card">
          <EmptyState art="search" title="No plans match" sub="Try a different search term." />
        </div>
      ) : (
        <div className="tile-grid" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))' }}>
          {list.map((p) => (
            <PlanCard key={p.id} plan={p} readOnly />
          ))}
        </div>
      )}
    </div>
  );
}
