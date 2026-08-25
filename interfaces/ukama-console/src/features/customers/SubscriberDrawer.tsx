/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import { useGetPackagesForSimQuery } from '@/client/graphql/packages.generated';
import { useToggleSimServiceStatusMutation } from '@/client/graphql/sims.generated';
import AppDrawer, { DetailRow } from '@/components/AppDrawer';
import AppModal from '@/components/AppModal';
import Meter from '@/components/Meter';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';
import { formatDate, parseTimestamp } from '@/lib/parsers';
import { bytesToGB, formatBytes } from '@/lib/usage';
import {
  NO_DATA_PLANS_MESSAGE,
  NO_POOL_SIMS_MESSAGE,
  useAvailablePoolSims,
  useDataPlans,
} from '@/lib/sim-pool';
import { useNetworkId } from '@/lib/useNetworkId';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import MoreVertRounded from '@mui/icons-material/MoreVertRounded';
import ReceiptLongRounded from '@mui/icons-material/ReceiptLongRounded';
import SimCardRounded from '@mui/icons-material/SimCardRounded';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import Skeleton from '@mui/material/Skeleton';
import { useMemo, useState } from 'react';
import AllocateSimDialog from './AllocateSimDialog';
import ReceiptDialog from './ReceiptDialog';
import TopUpDialog from './TopUpDialog';

type SimPackage = {
  id: string;
  package_id: string;
  start_date: string;
  end_date: string;
  is_currently_in_use: boolean;
};

type PackageKind = 'current' | 'upcoming' | 'ended';

const classify = (p: SimPackage, now: number): PackageKind => {
  if (p.is_currently_in_use) return 'current';
  // A queued package chains off the previous package's end date, so its start
  // can already be in the past while it is still waiting to activate. Decide
  // "ended" by the END date, not the start — otherwise a queued package whose
  // end is still in the future gets wrongly shown as Ended.
  const end = parseTimestamp(p.end_date);
  if (!Number.isNaN(end) && end > now) return 'upcoming';
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

function PackageMenu({ onViewReceipt }: { onViewReceipt: () => void }) {
  const [anchor, setAnchor] = useState<null | HTMLElement>(null);
  return (
    <>
      <IconButton
        size="small"
        aria-label="Package options"
        onClick={(e) => setAnchor(e.currentTarget)}
        sx={{ color: 'var(--uk-ink-3)' }}
      >
        <MoreVertRounded sx={{ fontSize: 20 }} />
      </IconButton>
      <Menu anchorEl={anchor} open={!!anchor} onClose={() => setAnchor(null)}>
        <MenuItem
          onClick={() => {
            setAnchor(null);
            onViewReceipt();
          }}
        >
          <ReceiptLongRounded sx={{ fontSize: 18, mr: 1, color: 'var(--uk-ink-3)' }} />
          View receipt
        </MenuItem>
      </Menu>
    </>
  );
}

export default function SubscriberDrawer({
  sub,
  onClose,
  agent,
  onChanged,
}: {
  sub: Subscriber;
  onClose: () => void;
  /** Customer (agent) lens only — enables the Top-up / SIM management footer. */
  agent?: boolean;
  onChanged?: () => void;
}) {
  const toast = useToast();
  const [showTopUp, setShowTopUp] = useState(false);
  const [showAllocate, setShowAllocate] = useState(false);
  const [receiptPkg, setReceiptPkg] = useState<{
    packageId: string;
    planName: string;
  } | null>(null);
  const hasSim = !!sub.simId;
  const simActive = sub.sim === 'active';
  // Suspended SIMs are managed elsewhere — the active/inactive toggle is
  // disabled for them here.
  const simSuspended = sub.sim === 'suspended';
  const [confirmSim, setConfirmSim] = useState(false);

  const [toggleSim, { loading: togglingSim }] =
    useToggleSimServiceStatusMutation({
      onCompleted: (d) => {
        setConfirmSim(false);
        const res = d.toggleSimServiceStatus;
        if (res.success) {
          // Activation/deactivation runs asynchronously on the backend, so the
          // toast reports that the process has started, not that it's applied.
          toast(
            `SIM ${simActive ? 'deactivation' : 'activation'} process initiated`,
          );
          onChanged?.(); // refresh the list + this subscriber's sim status
        } else {
          toast(res.message ?? "Couldn't update SIM status");
        }
      },
      onError: (e) => {
        setConfirmSim(false);
        toast(e.message || "Couldn't update SIM status");
      },
    });

  const doToggleSim = () => {
    if (!sub.simId) return;
    toggleSim({
      variables: {
        data: {
          sim_id: sub.simId,
          status: simActive ? 'service_off' : 'service_on',
        },
      },
    });
  };

  // Allocating needs an available pool SIM and a data plan — guide the user to
  // set up whichever is missing instead of opening an unusable dialog. The same
  // plans list also resolves package ids → names for the packages section below.
  const networkId = useNetworkId();
  const { available: poolSims } = useAvailablePoolSims();
  const { packages: plans, available: dataPlans } = useDataPlans(networkId);
  const openAllocate = () => {
    if (poolSims === 0) {
      toast(NO_POOL_SIMS_MESSAGE);
      return;
    }
    if (dataPlans === 0) {
      toast(NO_DATA_PLANS_MESSAGE);
      return;
    }
    setShowAllocate(true);
  };
  // usage is raw bytes (-1 when unknown/none — clamp so we never render "-1").
  const usageBytes = Math.max(0, sub.usage);
  const usageLabel = formatBytes(usageBytes);
  const pct = sub.cap
    ? Math.min(100, (bytesToGB(usageBytes) / sub.cap) * 100)
    : 50;
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
  const planNameById = useMemo(() => {
    const m = new Map<string, string>();
    for (const p of plans) m.set(p.uuid, p.name);
    return m;
  }, [plans]);

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
                {usageLabel} of {sub.cap} GB used this cycle
              </div>
            </>
          ) : (
            <div
              className="tnum"
              style={{ fontSize: 12.5, color: 'var(--uk-ink-2)' }}
            >
              {usageLabel} used · unlimited
            </div>
          )}
        </div>

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
                      <div style={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <StatusBadge status={badge.status}>
                          {badge.label}
                        </StatusBadge>
                        {hasSim && (
                          <PackageMenu
                            onViewReceipt={() =>
                              setReceiptPkg({
                                packageId: p.package_id,
                                planName:
                                  planNameById.get(p.package_id) ?? p.package_id,
                              })
                            }
                          />
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </div>

      {agent && (
        <div
          style={{
            padding: '14px 24px',
            borderTop: '1px solid var(--uk-line)',
            display: 'flex',
            gap: 10,
          }}
        >
          {hasSim ? (
            <>
              <Button
                variant="contained"
                startIcon={<AddCardRounded />}
                sx={{ flex: 1 }}
                disabled={!simActive}
                title={!simActive ? 'SIM must be active to top up' : undefined}
                onClick={() => setShowTopUp(true)}
              >
                Top up
              </Button>
              <Button
                variant="outlined"
                color={simActive ? 'error' : 'primary'}
                startIcon={<SimCardRounded />}
                sx={{ flex: 1 }}
                disabled={togglingSim || simSuspended}
                title={simSuspended ? 'SIM is suspended' : undefined}
                onClick={() => setConfirmSim(true)}
              >
                {simSuspended
                  ? 'SIM suspended'
                  : simActive
                    ? 'Deactivate SIM'
                    : 'Activate SIM'}
              </Button>
            </>
          ) : (
            <Button
              variant="contained"
              startIcon={<SimCardRounded />}
              sx={{ flex: 1 }}
              onClick={openAllocate}
            >
              Allocate a SIM
            </Button>
          )}
        </div>
      )}

      {confirmSim && (
        <AppModal
          title={simActive ? 'Deactivate SIM' : 'Activate SIM'}
          width={420}
          onClose={() => {
            if (!togglingSim) setConfirmSim(false);
          }}
          footer={
            <>
              <Button
                color="inherit"
                sx={{ color: 'var(--uk-ink-3)' }}
                disabled={togglingSim}
                onClick={() => setConfirmSim(false)}
              >
                Cancel
              </Button>
              <Button
                variant="contained"
                color={simActive ? 'error' : 'primary'}
                startIcon={<SimCardRounded />}
                disabled={togglingSim}
                onClick={doToggleSim}
              >
                {togglingSim
                  ? 'Updating…'
                  : simActive
                    ? 'Deactivate SIM'
                    : 'Activate SIM'}
              </Button>
            </>
          }
        >
          <div
            style={{ fontSize: 14, color: 'var(--uk-ink-2)', lineHeight: 1.55 }}
          >
            {simActive ? (
              <>
                This will deactivate the SIM for{' '}
                <b style={{ color: 'var(--uk-ink)' }}>{sub.name}</b>. They&apos;ll
                lose connectivity until the SIM is reactivated.
              </>
            ) : (
              <>
                This will activate the SIM for{' '}
                <b style={{ color: 'var(--uk-ink)' }}>{sub.name}</b> and restore
                their connectivity.
              </>
            )}
          </div>
        </AppModal>
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

      {receiptPkg && sub.simId && (
        <ReceiptDialog
          simId={sub.simId}
          packageId={receiptPkg.packageId}
          planName={receiptPkg.planName}
          onClose={() => setReceiptPkg(null)}
        />
      )}
    </AppDrawer>
  );
}
