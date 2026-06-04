/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Create data plan dialog (form-dialogs.jsx) — react-hook-form + zod. */
import Button from '@mui/material/Button';
import AddRounded from '@mui/icons-material/AddRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { createPlanSchema } from './schemas';
import type { CreatePlanValues } from './schemas';

export default function CreatePlanDialog({ onClose }: { onClose: () => void }) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<CreatePlanValues>({
    resolver: zodResolver(createPlanSchema),
    mode: 'onChange',
    defaultValues: { name: '', unit: 'GB', days: 30 } as Partial<CreatePlanValues> as CreatePlanValues,
  });

  const submit = handleSubmit((v) => {
    onClose();
    toast(`Plan “${v.name}” created!`);
  });

  return (
    <AppModal
      title="Create data plan"
      width={540}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Cancel
          </Button>
          <Button variant="contained" startIcon={<AddRounded />} disabled={!isValid} onClick={submit}>
            Create plan
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 18px', textWrap: 'pretty' }}>
        Create a custom data plan that can be assigned to customers.
      </p>
      <Field label="Data plan name" required error={errors.name?.message}>
        <TextInput placeholder="e.g. Standard" invalid={!!errors.name} {...register('name')} />
      </Field>
      <Field label="Price" required error={errors.price?.message}>
        <TextInput placeholder="0.00" prefix="$" type="number" invalid={!!errors.price} {...register('price')} />
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
              { value: 'Unlimited', label: 'Unlimited' },
            ]}
            {...register('unit')}
          />
        </Field>
      </div>
      <Field label="Validity (days)" error={errors.days?.message}>
        <TextInput placeholder="30" type="number" {...register('days')} />
      </Field>
    </AppModal>
  );
}
