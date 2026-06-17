/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Create / edit data plan dialog — react-hook-form + zod.
 * Create: addPackage with all fields. Edit: the gateway only allows changing
 * the name (UpdatePackageInputDto = { name, active }), so the other fields
 * are shown prefilled but read-only, and we call updatePackage.
 *
 * Plan names must be unique, so the name is validated live against
 * isPackageNameAvailable (debounced) and the submit button is gated on a
 * confirmed-available name. A plan is either org-wide ("Available within an
 * org" → no networkId) or scoped to a single network.
 */
import { useEffect, useState } from 'react';
import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import CircularProgress from '@mui/material/CircularProgress';
import FormControlLabel from '@mui/material/FormControlLabel';
import AddRounded from '@mui/icons-material/AddRounded';
import SaveRounded from '@mui/icons-material/SaveRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import {
  GetPackagesDocument,
  useAddPackageMutation,
  useIsPackageNameAvailableLazyQuery,
  useUpdatePacakgeMutation,
  type PackageFragment,
} from '@/client/graphql/packages.generated';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { useAuth } from '@/lib/auth/context';
import { useCurrency } from '@/lib/currency';
import { createPlanSchema, VALIDITY_OPTIONS } from './schemas';
import type { CreatePlanValues } from './schemas';

const normalizeUnit = (u: string): 'GB' | 'MB' =>
  u?.toUpperCase() === 'MB' ? 'MB' : 'GB';

const normalizeDays = (d: number): number =>
  [1, 7, 30].includes(Math.round(d)) ? Math.round(d) : 30;

const validityLabel = (days: number): string =>
  VALIDITY_OPTIONS.find((o) => o.value === String(days))?.label ??
  `${days} days`;

/** Live name-availability state for the inline validation line. */
type NameState = 'idle' | 'checking' | 'available' | 'taken';

export default function CreatePlanDialog({
  pkg,
  onClose,
}: {
  pkg?: PackageFragment | null;
  onClose: () => void;
}) {
  const toast = useToast();
  const user = useAuth();
  const { symbol } = useCurrency();
  const editing = Boolean(pkg);

  const { data: networksData, loading: networksLoading } = useGetNetworksQuery();
  const networkOptions =
    networksData?.getNetworks.networks.map((n) => ({
      value: n.id,
      label: n.name,
    })) ?? [];

  const {
    register,
    control,
    watch,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<CreatePlanValues>({
    resolver: zodResolver(createPlanSchema),
    mode: 'onChange',
    defaultValues: pkg
      ? {
          name: pkg.name,
          price: pkg.amount,
          data: pkg.dataVolume,
          unit: normalizeUnit(pkg.dataUnit),
          days: normalizeDays(pkg.duration),
          availableWithinOrg: !pkg.networkId,
          networkId: pkg.networkId ?? '',
        }
      : ({
          name: '',
          unit: 'GB',
          days: 30,
          availableWithinOrg: true,
          networkId: '',
        } as CreatePlanValues),
  });

  // eslint-disable-next-line react-hooks/incompatible-library -- RHF watch() is intentionally incompatible
  const orgWide = watch('availableWithinOrg');
  const nameValue = watch('name');
  const trimmedName = (nameValue ?? '').trim();
  const nameUnchanged = editing && pkg ? trimmedName === pkg.name : false;

  // ---- live name-availability check (debounced) ----
  const [nameState, setNameState] = useState<NameState>('idle');
  const [checkName] = useIsPackageNameAvailableLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (res) =>
      setNameState(res.isPackageNameAvailable.isAvailable ? 'available' : 'taken'),
    onError: () => setNameState('idle'),
  });

  useEffect(() => {
    // Empty name (zod handles "required") or an unchanged name in edit mode
    // never needs a round-trip.
    if (!trimmedName || nameUnchanged) {
      setNameState('idle');
      return;
    }
    setNameState('checking');
    const t = setTimeout(() => {
      void checkName({ variables: { name: trimmedName } });
    }, 400);
    return () => clearTimeout(t);
  }, [trimmedName, nameUnchanged, checkName]);

  const refetch = [{ query: GetPackagesDocument }];

  const [addPackage, { loading: adding }] = useAddPackageMutation({
    refetchQueries: refetch,
    onCompleted: (res) => {
      toast(`Plan “${res.addPackage.name}” created!`);
      onClose();
    },
    onError: (err) => toast(err.message),
  });

  const [updatePackage, { loading: updating }] = useUpdatePacakgeMutation({
    refetchQueries: refetch,
    onCompleted: (res) => {
      toast(`Plan “${res.updatePackage.name}” updated!`);
      onClose();
    },
    onError: (err) => toast(err.message),
  });

  const saving = adding || updating;

  // The name is acceptable when it's confirmed available, or (when editing)
  // left unchanged.
  const nameOk = nameUnchanged || nameState === 'available';

  const submit = handleSubmit((v) => {
    if (editing && pkg) {
      void updatePackage({
        variables: {
          packageId: pkg.uuid,
          data: { name: v.name, active: pkg.active },
        },
      });
      return;
    }
    if (!user?.currency) {
      toast('Your organization has no currency set. Please contact support.');
      return;
    }
    void addPackage({
      variables: {
        data: {
          name: v.name,
          amount: v.price,
          dataVolume: v.data,
          dataUnit: v.unit,
          duration: v.days,
          country: user.country ?? '',
          currency: user.currency,
          // Org-wide → no network; otherwise the selected network.
          networkId: v.availableWithinOrg ? '' : (v.networkId ?? ''),
        },
      },
    });
  });

  // Inline availability message under the name field (zod errors take priority).
  const availabilityLine =
    errors.name || !trimmedName || nameUnchanged ? null : (
      <div
        style={{
          marginTop: 5,
          fontSize: 11.5,
          color:
            nameState === 'taken'
              ? 'var(--uk-error)'
              : nameState === 'available'
                ? 'var(--uk-success)'
                : 'var(--uk-ink-3)',
        }}
      >
        {nameState === 'checking'
          ? 'Checking availability…'
          : nameState === 'available'
            ? '✓ Name is available'
            : nameState === 'taken'
              ? 'That plan name is already taken'
              : null}
      </div>
    );

  return (
    <AppModal
      title={editing ? 'Edit data plan' : 'Create data plan'}
      width={540}
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
              ) : editing ? (
                <SaveRounded />
              ) : (
                <AddRounded />
              )
            }
            disabled={!isValid || saving || !nameOk}
            onClick={submit}
          >
            {saving
              ? editing
                ? 'Saving…'
                : 'Creating…'
              : editing
                ? 'Save changes'
                : 'Create plan'}
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 18px', textWrap: 'pretty' }}>
        {editing
          ? 'Update the plan name. Price, data, and validity are fixed once a plan is created.'
          : 'Create a custom data plan that can be assigned to customers.'}
      </p>
      <Field label="Data plan name" required error={errors.name?.message}>
        <TextInput
          placeholder="e.g. Standard"
          invalid={!!errors.name || nameState === 'taken'}
          {...register('name')}
        />
        {availabilityLine}
      </Field>

      {editing && pkg ? (
        <>
          <Field label="Price">
            <div className="ff-readonly">
              {symbol}
              {pkg.amount}
            </div>
          </Field>
          <div className="ff-grid2">
            <Field label="Data volume">
              <div className="ff-readonly">{pkg.dataVolume}</div>
            </Field>
            <Field label="Unit">
              <div className="ff-readonly">{normalizeUnit(pkg.dataUnit)}</div>
            </Field>
          </div>
          <Field label="Validity">
            <div className="ff-readonly">
              {validityLabel(normalizeDays(pkg.duration))}
            </div>
          </Field>
        </>
      ) : (
        <>
          <Field label="Price" required error={errors.price?.message}>
            <TextInput placeholder="0.00" prefix={symbol} type="number" invalid={!!errors.price} {...register('price')} />
          </Field>
          <div className="ff-grid2">
            <Field label="Data volume" required error={errors.data?.message}>
              <TextInput placeholder="20" type="number" invalid={!!errors.data} {...register('data')} />
            </Field>
            <Field label="Unit">
              <SelectInput
                options={[
                  { value: 'GB', label: 'GB' },
                  { value: 'MB', label: 'MB' },
                ]}
                {...register('unit')}
              />
            </Field>
          </div>
          <Field label="Validity" required error={errors.days?.message}>
            <SelectInput
              options={VALIDITY_OPTIONS.map((o) => ({ value: o.value, label: o.label }))}
              {...register('days')}
            />
          </Field>

          <Controller
            control={control}
            name="availableWithinOrg"
            render={({ field }) => (
              <FormControlLabel
                sx={{ ml: 0, mb: '14px', alignSelf: 'baseline' }}
                control={
                  <Checkbox
                    checked={field.value}
                    onChange={(e) => field.onChange(e.target.checked)}
                    sx={{ p: 0, pr: 1.5 }}
                  />
                }
                label="Available within an org"
              />
            )}
          />

          {!orgWide && (
            <Field label="Network" required error={errors.networkId?.message}>
              <SelectInput
                placeholder={networksLoading ? 'Loading networks…' : 'Select a network'}
                options={networkOptions}
                invalid={!!errors.networkId}
                {...register('networkId')}
              />
            </Field>
          )}
        </>
      )}
    </AppModal>
  );
}
