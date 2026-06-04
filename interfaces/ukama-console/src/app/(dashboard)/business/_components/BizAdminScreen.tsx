/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Admin — utility launcher into the detailed management pages (biz-ops.jsx). */
import Link from 'next/link';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import InfoRounded from '@mui/icons-material/InfoRounded';
import PageHeader from '@/components/PageHeader';
import SectionCard from '@/components/SectionCard';
import { BIZ_ADMIN } from '@/data';

export default function BizAdminScreen() {
  return (
    <div className="page">
      <PageHeader title="Admin" sub="Detailed technical and resource management." />
      <div
        className="card card-pad"
        style={{
          marginBottom: 'var(--uk-gap)',
          display: 'flex',
          gap: 13,
          alignItems: 'flex-start',
        }}
      >
        <InfoRounded sx={{ color: 'var(--uk-ac)', fontSize: 22, flex: 'none', mt: '1px' }} />
        <div>
          <div style={{ fontSize: 14, fontWeight: 600 }}>Admin area</div>
          <div
            style={{ fontSize: 13, color: 'var(--uk-ink-2)', marginTop: 3, textWrap: 'pretty' }}
          >
            The detailed technical and resource pages live here. Admin is the utility area —
            not the default operator experience.
          </div>
        </div>
      </div>
      <div className="tile-grid" style={{ gridTemplateColumns: '1fr 1fr' }}>
        {BIZ_ADMIN.map((sec) => (
          <SectionCard key={sec.group} title={sec.group}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              {sec.items.map((it) => (
                <Link
                  key={it.name}
                  href={it.href}
                  className="admin-launch"
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 12,
                    textAlign: 'left',
                    padding: '13px 15px',
                    borderRadius: 10,
                    cursor: 'pointer',
                    border: '1px solid var(--uk-line)',
                    background: 'var(--uk-page)',
                    textDecoration: 'none',
                    transition: '.12s',
                  }}
                >
                  <div style={{ flex: 1 }}>
                    <div style={{ fontSize: 13.5, fontWeight: 600, color: 'var(--uk-ink)' }}>
                      {it.name}
                    </div>
                    <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)', marginTop: 1 }}>
                      {it.desc}
                    </div>
                  </div>
                  <ChevronRightRounded sx={{ fontSize: 20, color: 'var(--uk-ink-3)' }} />
                </Link>
              ))}
            </div>
          </SectionCard>
        ))}
      </div>
    </div>
  );
}
