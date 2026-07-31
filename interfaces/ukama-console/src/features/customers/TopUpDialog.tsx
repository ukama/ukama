/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import { useAddPaymentMutation } from '@/client/graphql/packages.generated';
import AppModal from '@/components/AppModal';
import { Field, SelectInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import type { Subscriber } from '@/data';
import { useCurrency } from '@/lib/currency';
import { durationUnitLabel } from '@/lib/duration';
import { useDataPlans } from '@/lib/sim-pool';
import { useNetworkId } from '@/lib/useNetworkId';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useMemo, useState } from 'react';

export default function TopUpDialog({
  sub,
  onClose,
  onDone,
}: {
  sub: Subscriber;
  onClose: () => void;
  onDone?: () => void;
}) {
  const toast = useToast();
  const { symbol, code } = useCurrency();
  const networkId = useNetworkId();
  const [planId, setPlanId] = useState('');

  const { packages, loading } = useDataPlans(networkId);
  const planOptions = useMemo(
    () =>
      packages.map((p) => ({
        value: p.uuid,
        label: `${p.name} · ${symbol}${p.amount}/${durationUnitLabel(p.duration)}`,
      })),
    [packages, symbol],
  );

  // Record a cash payment for the package. The payments system settles cash
  // immediately and emits payment.success, which sim-manager consumes to
  // allocate the package to the subscriber's SIM (no direct addPackagesToSim).
  const [addPayment, { loading: saving }] = useAddPaymentMutation({
    onCompleted: () => {
      toast(`Topped up ${sub.name}`);
      onDone?.();
      onClose();
    },
    onError: (err) => toast(err.message),
  });

  const submit = () => {
    if (!sub.simId) {
      toast('This customer has no SIM to top up.');
      return;
    }
    if (!planId) {
      toast('Select a data plan.');
      return;
    }
    const plan = packages.find((p) => p.uuid === planId);
    if (!plan) {
      toast('Select a valid data plan.');
      return;
    }
    void addPayment({
      variables: {
        data: {
          itemId: planId,
          itemType: 'package',
          paymentMethod: 'cash',
          amount: String(plan.amount),
          currency: code,
          payerEmail: sub.email,
          payerName: sub.name,
          payerPhone: sub.phone,
          description: plan.name,
          country: plan.country,
          sim: sub.simId,
        },
      },
    });
  };

  return (
    <AppModal
      title="Top up data"
      width={480}
      onClose={onClose}
      footer={
        <>
          <Button
            color="inherit"
            sx={{ color: 'var(--uk-ink-3)' }}
            onClick={onClose}
            disabled={saving}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            startIcon={
              saving ? (
                <CircularProgress size={16} color="inherit" />
              ) : (
                <AddCardRounded />
              )
            }
            disabled={!planId || saving || loading}
            onClick={submit}
          >
            {saving ? 'Topping up…' : 'Top up'}
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 18px', textWrap: 'pretty' }}>
        Add a data plan to <strong>{sub.name}</strong>&apos;s SIM. It activates
        shortly after you confirm.
      </p>
      <Field label="Data plan" required>
        <SelectInput
          placeholder={loading ? 'Loading plans…' : 'Select a plan'}
          options={planOptions}
          onChange={(e) => setPlanId(e.target.value)}
        />
      </Field>
    </AppModal>
  );
}
