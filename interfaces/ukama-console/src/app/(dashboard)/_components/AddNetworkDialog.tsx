/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Add-network dialog (top-bar switcher). Names + creates an additional network
 * inline — same validation as the /configure network step — then selects it.
 * The new network is never made default (isDefault: false) so it can't steal
 * the org's default from the switcher.
 */
import { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';

import {
  GetNetworksDocument,
  useAddNetworkMutation,
} from '@/client/graphql/networks.generated';
import { OnboardingStatusDocument } from '@/client/graphql/onboarding-status.generated';
import AppModal from '@/components/AppModal';
import { Field, TextInput } from '@/components/form/FormField';
import { useToast } from '@/components/ToastProvider';
import { useUiPrefs } from '@/lib/store';

const schema = z.object({
  name: z
    .string()
    .min(3, 'At least 3 characters')
    .max(40, 'At most 40 characters')
    .regex(
      /^[a-z0-9-]+$/,
      'Lowercase letters, numbers, and "-" only (no spaces).',
    ),
});

type FormValues = z.infer<typeof schema>;

export default function AddNetworkDialog({
  onClose,
  onCreated,
}: {
  onClose: () => void;
  onCreated?: (id: string) => void;
}) {
  const toast = useToast();
  const setNetworkId = useUiPrefs((s) => s.setNetworkId);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  const [addNetwork, { loading }] = useAddNetworkMutation({
    refetchQueries: [
      { query: OnboardingStatusDocument },
      { query: GetNetworksDocument },
    ],
    onCompleted: (res) => {
      setNetworkId(res.addNetwork.id);
      toast(`Network "${res.addNetwork.name}" created.`);
      onCreated?.(res.addNetwork.id);
      onClose();
    },
    onError: (err) => setSubmitError(err.message),
  });

  const onSubmit = (values: FormValues) => {
    setSubmitError(null);
    // An additional network never steals the default.
    void addNetwork({
      variables: { data: { name: values.name, isDefault: false } },
    });
  };

  return (
    <AppModal
      title="Add network"
      width={460}
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
              ) : undefined
            }
            disabled={loading}
            onClick={() => void handleSubmit(onSubmit)()}
          >
            {loading ? 'Creating…' : 'Create network'}
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
        A network is made up of one or more sites of Ukama hardware, allowing
        you to connect to the cellular internet.
      </p>
      <form onSubmit={(e) => void handleSubmit(onSubmit)(e)}>
        <Field label="Network name" required error={errors.name?.message}>
          <TextInput
            placeholder="network-name"
            invalid={Boolean(errors.name)}
            autoFocus
            {...register('name')}
          />
        </Field>
        {submitError && <p className="cfg-error">{submitError}</p>}
      </form>
    </AppModal>
  );
}
