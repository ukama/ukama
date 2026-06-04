/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Add customer dialog (form-dialogs.jsx) — react-hook-form + zod. */
import Button from '@mui/material/Button';
import PersonAddRounded from '@mui/icons-material/PersonAddRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { PLANS } from '@/data';
import { addCustomerSchema } from './schemas';
import type { AddCustomerValues } from './schemas';

export default function AddCustomerDialog({ onClose }: { onClose: () => void }) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<AddCustomerValues>({
    resolver: zodResolver(addCustomerSchema),
    mode: 'onChange',
    defaultValues: { first: '', last: '', mobile: '', email: '', planId: '', sim: '' },
  });

  const submit = handleSubmit((v) => {
    onClose();
    toast(`${v.first}${v.last ? ' ' + v.last : ''} added successfully!`);
  });

  return (
    <AppModal
      title="Add customer"
      width={560}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Cancel
          </Button>
          <Button
            variant="contained"
            startIcon={<PersonAddRounded />}
            disabled={!isValid}
            onClick={submit}
          >
            Add customer
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 18px', textWrap: 'pretty' }}>
        Add a customer to your network and assign them a SIM and data plan.
      </p>
      <div className="ff-grid2">
        <Field label="First name" required error={errors.first?.message}>
          <TextInput placeholder="John" invalid={!!errors.first} {...register('first')} />
        </Field>
        <Field label="Last name">
          <TextInput placeholder="Doe" {...register('last')} />
        </Field>
      </div>
      <div className="ff-grid2">
        <Field label="Mobile number" required error={errors.mobile?.message}>
          <TextInput placeholder="+260 97 000 0000" invalid={!!errors.mobile} {...register('mobile')} />
        </Field>
        <Field label="Email" error={errors.email?.message}>
          <TextInput placeholder="name@email.com" type="email" invalid={!!errors.email} {...register('email')} />
        </Field>
      </div>
      <Field label="Data plan">
        <SelectInput
          placeholder="Select a plan"
          options={PLANS.map((p) => ({
            value: p.id,
            label: `${p.name} · $${p.price}/${p.days === 1 ? 'day' : 'mo'}`,
          }))}
          {...register('planId')}
        />
      </Field>
      <Field label="SIM ICCID" hint="Leave blank to assign automatically from the pool">
        <TextInput placeholder="8926 0010 0000 0000" {...register('sim')} />
      </Field>
    </AppModal>
  );
}
