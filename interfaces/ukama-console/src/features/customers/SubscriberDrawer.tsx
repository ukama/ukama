/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import {
  useGetPackagesForSimQuery,
  useGetPackagesQuery,
} from '@/client/graphql/packages.generated';
import AppDrawer, { DetailRow } from '@/components/AppDrawer';
import Meter from '@/components/Meter';
import StatusBadge from '@/components/StatusBadge';
import type { Subscriber } from '@/data';
import { formatDate, parseTimestamp } from '@/lib/parsers';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import SimCardRounded from '@mui/icons-material/SimCardRounded';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useMemo, useState } from 'react';
import AllocateSimDialog from './AllocateSimDialog';
import TopUpDialog from './TopUpDialog';

type SimPackage = {
  id: string;
  package_id: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
};

type PackageKind = 'current' | 'upcoming' | 'ended';

const classify = (p: SimPackage, now: number): PackageKind => {
  if (p.is_active) return 'current';
  const start = parseTimestamp(p.start_date);
  if (!Number.isNaN(start) && start > now) return 'upcoming';
  return 'ended';
};

const KIND_BADGE: Record<PackageKind, { status: string; label: string }> = {
  current: { status: 'active', label: 'Current' },
  upcoming: { status: 'pending', label: 'Upcoming' },
  ended: { status: 'inactive', label: 'Ended' },
};

const KIND_ORDER: Record<PackageKind, number> = {
  current: 0,
  upcoming: 1,
  ended: 2,
};

export default function SubscriberDrawer({
  sub,
  onClose,
  readOnly,
  onChanged,
}: {
  sub: Subscriber;
  onClose: () => void;
  readOnly?: boolean;
  onChanged?: () => void;
}) {
  const [showTopUp, setShowTopUp] = useState(false);
  const [showAllocate, setShowAllocate] = useState(false);
  const hasSim = !!sub.simId;
  const pct = sub.cap ? Math.min(100, (sub.usage / sub.cap) * 100) : 50;
  const initials = sub.name
    .split(' ')
    .map((x) => x[0])
    .join('');

  // Packages on this subscriber's SIM, with plan names resolved.
  const {
    data: simPkgData,
    loading: pkgLoading,
    refetch: refetchPkgs,
  } = useGetPackagesForSimQuery({
    variables: { data: { sim_id: sub.simId ?? '' } },
    skip: !sub.simId,
    fetchPolicy: 'cache-and-network',
  });
  const { data: pkgData } = useGetPackagesQuery();
  const planNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const p of pkgData?.getPackages.packages ?? []) m.set(p.uuid, p.name);
    return m;
  }, [pkgData]);

  const [now] = useState(() => Date.now());
  const packages = [...(simPkgData?.getPackagesForSim.packages ?? [])]
    .map((p) => ({ ...p, kind: classify(p, now) }))
    .sort(
      (a, b) =>
        KIND_ORDER[a.kind] - KIND_ORDER[b.kind] ||
        parseTimestamp(a.start_date) - parseTimestamp(b.start_date),
    );

  return (
    <AppDrawer onClose={onClose} width={430}>
      <div
        style={{
          padding: '20px 24px 16px',
          borderBottom: '1px solid var(--uk-line)',
        }}
      >
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
          }}
        >
          <div style={{ display: 'flex', gap: 13, alignItems: 'center' }}>
            <span
              className="av-sm"
              style={{ width: 46, height: 46, fontSize: 16 }}
            >
              {initials}
            </span>
            <div>
              <div
                style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 18,
                  fontWeight: 500,
                }}
              >
                {sub.name}
              </div>
              <div
                className="tnum"
                style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}
              >
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
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              marginBottom: 8,
            }}
          >
            <span style={{ fontSize: 13, fontWeight: 600 }}>{sub.plan}</span>
            <StatusBadge status={sub.sim === 'suspended' ? 'pending' : sub.sim}>
              {sub.sim === 'suspended' ? 'Suspended' : undefined}
            </StatusBadge>
          </div>
          {sub.cap ? (
            <>
              <Meter
                value={pct}
                color={pct > 90 ? 'var(--uk-orange)' : undefined}
              />
              <div
                className="tnum"
                style={{
                  fontSize: 12.5,
                  color: 'var(--uk-ink-2)',
                  marginTop: 7,
                }}
              >
                {sub.usage} of {sub.cap} GB used this cycle
              </div>
            </>
          ) : (
            <div
              className="tnum"
              style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}
            >
              {sub.usage} GB used · unlimited
            </div>
          )}
        </div>

        <DetailRow k="Site" v={sub.site} />
        <DetailRow k="ICCID" v={sub.iccid} />
        <DetailRow
          k="SIM status"
          v={<span style={{ textTransform: 'capitalize' }}>{sub.sim}</span>}
        />
        <DetailRow k="Phone" v={sub.phone} />

        {sub.simId && (
          <div style={{ marginTop: 18 }}>
            <div
              style={{
                fontSize: 12,
                fontWeight: 600,
                color: 'var(--uk-ink-3)',
                textTransform: 'uppercase',
                letterSpacing: '0.06em',
                marginBottom: 8,
              }}
            >
              Packages
            </div>
            {pkgLoading && packages.length === 0 ? (
              <Skeleton variant="rounded" height={60} />
            ) : packages.length === 0 ? (
              <div style={{ fontSize: 13, color: 'var(--uk-ink-3)' }}>
                No packages assigned yet.
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                {packages.map((p) => {
                  const badge = KIND_BADGE[p.kind];
                  return (
                    <div
                      key={p.id}
                      className="card card-pad"
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        gap: 10,
                        padding: '10px 12px',
                      }}
                    >
                      <div style={{ minWidth: 0 }}>
                        <div style={{ fontSize: 13.5, fontWeight: 600 }}>
                          {planNameById.get(p.package_id) ?? p.package_id}
                        </div>
                        <div
                          className="tnum"
                          style={{
                            fontSize: 12,
                            color: 'var(--uk-ink-3)',
                            marginTop: 2,
                          }}
                        >
                          {formatDate(p.start_date)} – {formatDate(p.end_date)}
                        </div>
                      </div>
                      <StatusBadge status={badge.status}>
                        {badge.label}
                      </StatusBadge>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}
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
          {hasSim ? (
            <Button
              variant="contained"
              startIcon={<AddCardRounded />}
              sx={{ flex: 1 }}
              onClick={() => setShowTopUp(true)}
            >
              Top up
            </Button>
          ) : (
            <Button
              variant="contained"
              startIcon={<SimCardRounded />}
              sx={{ flex: 1 }}
              onClick={() => setShowAllocate(true)}
            >
              Allocate a SIM
            </Button>
          )}
        </div>
      )}

      {showTopUp && (
        <TopUpDialog
          sub={sub}
          onClose={() => setShowTopUp(false)}
          onDone={() => {
            void refetchPkgs();
            onChanged?.();
          }}
        />
      )}

      {showAllocate && (
        <AllocateSimDialog
          sub={sub}
          onClose={() => setShowAllocate(false)}
          onDone={() => {
            void refetchPkgs();
            onChanged?.();
          }}
        />
      )}
    </AppDrawer>
  );
}
