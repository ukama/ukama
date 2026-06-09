/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Invite member dialog — react-hook-form + zod, wired to createInvitation. */
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import SendRounded from '@mui/icons-material/SendRounded';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import {
  GetInvitationsDocument,
  useCreateInvitationMutation,
} from '@/client/graphql/invitations.generated';
import { TeamListDocument } from '@/client/graphql/team.generated';
import { Role_Type } from '@/client/graphql/types';
import AppModal from '@/components/AppModal';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { inviteMemberSchema } from '../_schemas/inviteMember';
import type { InviteMemberValues } from '../_schemas/inviteMember';

/** Friendly role labels → gateway Role_Type enum. */
const ROLE_OPTIONS: { label: InviteMemberValues['role']; value: Role_Type }[] = [
  { label: 'Owner', value: Role_Type.RoleOwner },
  { label: 'Administrator', value: Role_Type.RoleAdmin },
  { label: 'Network owner', value: Role_Type.RoleNetworkOwner },
  { label: 'Vendor', value: Role_Type.RoleVendor },
];

const roleEnum = (label: InviteMemberValues['role']): Role_Type =>
  ROLE_OPTIONS.find((o) => o.label === label)?.value ?? Role_Type.RoleUser;

/** Derive a placeholder name from the email's local part (user sets it later). */
const nameFromEmail = (email: string): string => email.split('@')[0] ?? email;

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

  const [createInvitation, { loading }] = useCreateInvitationMutation({
    // Refresh the team list (composite) + the raw invitations cache.
    refetchQueries: [
      { query: TeamListDocument },
      { query: GetInvitationsDocument },
    ],
    onCompleted: (res) => {
      toast(`Invitation sent to ${res.createInvitation.email}`);
      onClose();
    },
    onError: (err) => toast(err.message),
  });

  const submit = handleSubmit((v) => {
    void createInvitation({
      variables: {
        data: {
          email: v.email.toLowerCase(),
          name: nameFromEmail(v.email),
          role: roleEnum(v.role),
        },
      },
    });
  });

  return (
    <AppModal
      title="Invite member"
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
                <SendRounded />
              )
            }
            disabled={!isValid || loading}
            onClick={submit}
          >
            {loading ? 'Sending…' : 'Invite member'}
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
          options={ROLE_OPTIONS.map((r) => ({ value: r.label, label: r.label }))}
          {...register('role')}
        />
      </Field>
    </AppModal>
  );
}
