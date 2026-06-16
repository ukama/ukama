/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Allocate a SIM + data plan to an existing subscriber who has none yet.
 * Mirrors the allocation half of AddCustomerDialog: data plans come from
 * getPackages, the SIM dropdown lists unallocated pool SIMs (blank =
 * auto-assign from the pool).
 */
import { useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import SimCardRounded from '@mui/icons-material/SimCardRounded';

import { useGetPackagesQuery } from '@/client/graphql/packages.generated';
import {
  useAllocateSimMutation,
  useGetSimsFromPoolQuery,
} from '@/client/graphql/sims.generated';
import { Sim_Status, type Sim_Types } from '@/client/graphql/types';
import AppModal from '@/components/AppModal';
import { Field, SelectInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';
import { useCurrency } from '@/lib/currency';
import { publicEnv } from '@/lib/runtime-env';
import { useUiPrefs } from '@/lib/store';

export default function AllocateSimDialog({
  sub,
  onClose,
  onDone,
}: {
  sub: Subscriber;
  onClose: () => void;
  onDone?: () => void;
}) {
  const toast = useToast();
  const { symbol } = useCurrency();
  const networkId = useUiPrefs((s) => s.networkId);
  const [planId, setPlanId] = useState('');
  const [iccid, setIccid] = useState('');

  const { data: pkgData } = useGetPackagesQuery();
  const planOptions = useMemo(
    () =>
      (pkgData?.getPackages.packages ?? []).map((p) => ({
        value: p.uuid,
        label: `${p.name} · ${symbol}${p.amount}/${p.duration === 1 ? 'day' : 'mo'}`,
      })),
    [pkgData, symbol],
  );

  const { data: simData } = useGetSimsFromPoolQuery({
    variables: {
      data: { type: publicEnv().simType as Sim_Types, status: Sim_Status.All },
    },
  });
  const simOptions = useMemo(
    () => [
      { value: '', label: 'Auto-assign from pool' },
      ...(simData?.getSimsFromPool.sims ?? [])
        .filter((s) => !s.isAllocated && !s.isFailed)
        .map((s) => ({ value: s.iccid, label: s.iccid })),
    ],
    [simData],
  );

  const [allocateSim, { loading }] = useAllocateSimMutation();

  const submit = async () => {
    if (!networkId) {
      toast('No network selected. Please select a network first.');
      return;
    }
    if (!planId) {
      toast('Select a data plan to allocate a SIM.');
      return;
    }
    try {
      await allocateSim({
        variables: {
          data: {
            subscriber_id: sub.id,
            network_id: networkId,
            package_id: planId,
            sim_type: publicEnv().simType,
            iccid: iccid || undefined,
            traffic_policy: 0,
          },
        },
      });
      toast(`SIM allocated to ${sub.name}.`);
      onDone?.();
      onClose();
    } catch (e) {
      toast(e instanceof Error ? e.message : 'Could not allocate SIM');
    }
  };

  return (
    <AppModal
      title="Allocate a SIM"
      width={520}
      onClose={onClose}
      footer={
        <>
          <Button
            color="inherit"
            sx={{ color: 'var(--uk-ink-3)' }}
            onClick={onClose}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            startIcon={
              loading ? (
                <CircularProgress size={16} color="inherit" />
              ) : (
                <SimCardRounded />
              )
            }
            disabled={!planId || loading}
            onClick={() => void submit()}
          >
            {loading ? 'Allocating…' : 'Allocate SIM'}
          </Button>
        </>
      }
    >
      <p
        style={{
          fontSize: 13.5,
          color: 'var(--uk-ink-2)',
          lineHeight: 1.6,
          margin: '0 0 18px',
          textWrap: 'pretty',
        }}
      >
        Assign a SIM and data plan to{' '}
        <b style={{ color: 'var(--uk-ink)' }}>{sub.name}</b>.
      </p>
      <Field label="Data plan" required>
        <SelectInput
          placeholder="Select a plan"
          options={planOptions}
          onChange={(e) => setPlanId(e.target.value)}
        />
      </Field>
      <Field
        label="SIM"
        hint="Leave on auto-assign to take the next SIM from the pool"
      >
        <SelectInput
          options={simOptions}
          onChange={(e) => setIccid(e.target.value)}
        />
      </Field>
    </AppModal>
  );
}
