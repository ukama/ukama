/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Select-network step — used when a network already exists but the flow has no
 * networkid (e.g. configuring a node from the pool, or resuming setup without a
 * network in the URL). Lets the user pick which network to install the site
 * into, then forwards to the install step. Self-guards: with no networks it
 * forwards to the creation step instead.
 */
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Skeleton from '@mui/material/Skeleton';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';

import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { useUiPrefs } from '@/lib/store';
import ConfigureActions from '../_components/ConfigureActions';
import { stepUrl, useConfigureParams } from '../_components/state';

export default function ConfigureSelectNetworkPage() {
  const router = useRouter();
  const { flow, networkid } = useConfigureParams();
  const setNetworkId = useUiPrefs((s) => s.setNetworkId);

  // network-only: a stale cached empty list must not skip the (now populated)
  // selection or wrongly forward to creation.
  const { data, loading } = useGetNetworksQuery({
    fetchPolicy: 'network-only',
  });
  const networks = data?.getNetworks.networks ?? [];
  const fallback = networks.find((n) => n.isDefault) ?? networks[0];

  // The seeded selection (a networkid from the URL wins, else the org default)
  // is derived; `override` only holds an explicit user choice.
  const [override, setOverride] = useState<string | null>(null);
  const selected = override ?? (networkid || fallback?.id || '');

  // Self-guard: nothing to select → go create the first network.
  useEffect(() => {
    if (!loading && data && networks.length === 0) {
      router.replace(stepUrl('/configure/network', { flow }));
    }
  }, [loading, data, networks.length, flow, router]);

  const onNext = () => {
    if (!selected) return;
    setNetworkId(selected);
    router.push(stepUrl('/configure/install', { flow, networkid: selected }));
  };

  if (loading || networks.length === 0) {
    return (
      <>
        <Skeleton width="60%" height={42} />
        <Skeleton width="100%" height={28} />
        <Skeleton width="100%" height={64} sx={{ mt: 2 }} />
        <Skeleton width="100%" height={64} />
      </>
    );
  }

  return (
    <>
      <h1 className="cfg-title">Select a network</h1>
      <p className="cfg-copy">
        Choose the network you want to install this site into. You can add more
        sites to it later.
      </p>
      <div
        className="cfg-fields"
        style={{ display: 'flex', flexDirection: 'column', flex: 1 }}
      >
        <div
          role="radiogroup"
          aria-label="Network"
          style={{ display: 'flex', flexDirection: 'column', gap: 10 }}
        >
          {networks.map((n) => {
            const active = selected === n.id;
            return (
              <button
                key={n.id}
                type="button"
                role="radio"
                aria-checked={active}
                onClick={() => setOverride(n.id)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  gap: 12,
                  textAlign: 'left',
                  cursor: 'pointer',
                  padding: '14px 16px',
                  borderRadius: 'var(--uk-r-sm)',
                  border: `1px solid ${active ? 'var(--uk-ac)' : 'var(--uk-line)'}`,
                  background: active ? 'var(--uk-ac-soft)' : 'var(--uk-hover)',
                  color: 'var(--uk-ink)',
                  font: 'inherit',
                  outline: active ? '1px solid var(--uk-ac)' : 'none',
                }}
              >
                <span style={{ minWidth: 0 }}>
                  <span
                    style={{
                      display: 'block',
                      fontSize: 14.5,
                      fontWeight: 600,
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {n.name}
                  </span>
                  {n.isDefault && (
                    <span style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>
                      Default network
                    </span>
                  )}
                </span>
                {active && (
                  <CheckCircleRounded
                    sx={{ fontSize: 20, color: 'var(--uk-ac)', flex: 'none' }}
                  />
                )}
              </button>
            );
          })}
        </div>
        <ConfigureActions
          nextLabel="Continue"
          onNext={onNext}
          nextDisabled={!selected}
        />
      </div>
    </>
  );
}
