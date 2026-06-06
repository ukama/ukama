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
 */
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import AddRounded from '@mui/icons-material/AddRounded';
import SaveRounded from '@mui/icons-material/SaveRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  GetPackagesDocument,
  useAddPackageMutation,
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

  const {
    register,
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
        }
      : ({ name: '', unit: 'GB', days: 30 } as CreatePlanValues),
  });

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
        },
      },
    });
  });

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
            disabled={!isValid || saving}
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
        <TextInput placeholder="e.g. Standard" invalid={!!errors.name} {...register('name')} />
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
        </>
      )}
    </AppModal>
  );
}
