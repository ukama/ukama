/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Invite member dialog (form-dialogs.jsx) — react-hook-form + zod. */
import Button from '@mui/material/Button';
import SendRounded from '@mui/icons-material/SendRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { inviteMemberSchema } from '../_schemas/inviteMember';
import type { InviteMemberValues } from '../_schemas/inviteMember';

const ROLES = ['Owner', 'Administrator', 'Vendor', 'Network owner'] as const;

export default function InviteMemberDialog({ onClose }: { onClose: () => void }) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<InviteMemberValues>({
    resolver: zodResolver(inviteMemberSchema),
    mode: 'onChange',
  });

  const submit = handleSubmit((v) => {
    onClose();
    toast(`Invitation sent to ${v.email}`);
  });

  return (
    <AppModal
      title="Invite member"
      width={520}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Cancel
          </Button>
          <Button variant="contained" startIcon={<SendRounded />} disabled={!isValid} onClick={submit}>
            Invite member
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 18px', textWrap: 'pretty' }}>
        Invite a teammate to help operate this network. They’ll receive an email to set up
        their account.
      </p>
      <Field label="Email" required error={errors.email?.message}>
        <TextInput placeholder="name@email.com" type="email" invalid={!!errors.email} {...register('email')} />
      </Field>
      <Field label="Role" required error={errors.role?.message}>
        <SelectInput
          placeholder="Select a role"
          invalid={!!errors.role}
          options={ROLES.map((r) => ({ value: r, label: r }))}
          {...register('role')}
        />
      </Field>
    </AppModal>
  );
}
