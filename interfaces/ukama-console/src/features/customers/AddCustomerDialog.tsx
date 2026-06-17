/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Add customer dialog — creates a subscriber, then allocates a SIM with the
 * chosen data plan. Data plans come from getPackages; the SIM dropdown lists
 * unallocated pool SIMs (blank = auto-assign from the pool).
 */
import { useMemo } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useGetPackagesQuery } from '@/client/graphql/packages.generated';
import { useGetSimsFromPoolQuery } from '@/client/graphql/sims.generated';
import { useAddSubscriberMutation } from '@/client/graphql/subscribers.generated';
import { useAllocateSimMutation } from '@/client/graphql/sims.generated';
import { Sim_Status, type Sim_Types } from '@/client/graphql/types';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { useCurrency } from '@/lib/currency';
import { publicEnv } from '@/lib/runtime-env';
import { useUiPrefs } from '@/lib/store';
import { addCustomerSchema } from './schemas';
import type { AddCustomerValues } from './schemas';

export default function AddCustomerDialog({
  onClose,
  onAdded,
}: {
  onClose: () => void;
  onAdded?: () => void;
}) {
  const toast = useToast();
  const { symbol } = useCurrency();
  const networkId = useUiPrefs((s) => s.networkId);

  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<AddCustomerValues>({
    resolver: zodResolver(addCustomerSchema),
    mode: 'onChange',
    defaultValues: { first: '', last: '', email: '', planId: '', sim: '' },
  });

  // Data plans + unallocated pool SIMs for the dropdowns. Only this network's
  // plans + org-wide plans (the BFF filters when networkId is passed).
  const { data: pkgData } = useGetPackagesQuery({ variables: { networkId } });
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

  const [addSubscriber, { loading: adding }] = useAddSubscriberMutation();
  const [allocateSim, { loading: allocating }] = useAllocateSimMutation();
  const saving = adding || allocating;

  const submit = handleSubmit(async (v) => {
    if (!networkId) {
      toast('No network selected. Please select a network first.');
      return;
    }
    const name = `${v.first}${v.last ? ' ' + v.last : ''}`;
    try {
      const res = await addSubscriber({
        variables: {
          data: { name, email: v.email ?? '', network_id: networkId },
        },
      });
      const sub = res.data?.addSubscriber;
      if (!sub) throw new Error('Could not add customer');
      // Assign a SIM + data plan when a plan was selected.
      if (v.planId) {
        await allocateSim({
          variables: {
            data: {
              subscriber_id: sub.uuid,
              network_id: networkId,
              package_id: v.planId,
              sim_type: publicEnv().simType,
              iccid: v.sim || undefined,
              traffic_policy: 0,
            },
          },
        });
      }
      toast(`${name} added successfully!`);
      onAdded?.();
      onClose();
    } catch (e) {
      toast(e instanceof Error ? e.message : 'Could not add customer');
    }
  });

  return (
    <AppModal
      title="Add customer"
      width={560}
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
                <PersonAddRounded />
              )
            }
            disabled={!isValid || saving}
            onClick={() => void submit()}
          >
            {saving ? 'Adding…' : 'Add customer'}
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
        Add a customer to your network and assign them a SIM and data plan.
      </p>
      <div className="ff-grid2">
        <Field label="First name" required error={errors.first?.message}>
          <TextInput
            placeholder="John"
            invalid={!!errors.first}
            {...register('first')}
          />
        </Field>
        <Field label="Last name">
          <TextInput placeholder="Doe" {...register('last')} />
        </Field>
      </div>
      <Field label="Email" error={errors.email?.message}>
        <TextInput
          placeholder="name@email.com"
          type="email"
          invalid={!!errors.email}
          {...register('email')}
        />
      </Field>
      <Field label="Data plan">
        <SelectInput
          placeholder="Select a plan"
          options={planOptions}
          {...register('planId')}
        />
      </Field>
      <Field
        label="SIM"
        hint="Leave on auto-assign to take the next SIM from the pool"
      >
        <SelectInput options={simOptions} {...register('sim')} />
      </Field>
    </AppModal>
  );
}
