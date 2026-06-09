/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Network creation — an independent, SELF-GUARDING step (legacy parity):
 * if a network already exists it never renders, it forwards to the install
 * step. A stale saved resume URL pointing here therefore auto-corrects.
 */
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import Skeleton from '@mui/material/Skeleton';

import {
  GetNetworksDocument,
  useAddNetworkMutation,
  useGetNetworksQuery,
} from '@/client/graphql/networks.generated';
import { OnboardingStatusDocument } from '@/client/graphql/onboarding-status.generated';
import { Field, TextInput } from '@/components/form/FormField';
import { useUiPrefs } from '@/lib/store';
import ConfigureActions from '../_components/ConfigureActions';
import { stepUrl, useConfigureParams } from '../_components/state';

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

export default function ConfigureNetworkPage() {
  const router = useRouter();
  const { flow } = useConfigureParams();
  const setNetworkId = useUiPrefs((s) => s.setNetworkId);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({ resolver: zodResolver(schema) });

  // Self-guard (onboarding only): a network already exists → this step is
  // done, forward. flow=add-network intentionally creates additional
  // networks, so the guard must not apply there. network-only: a stale
  // cached empty list must not re-show creation.
  const isAddFlow = flow === 'add-network';
  const { data, loading } = useGetNetworksQuery({
    fetchPolicy: 'network-only',
  });
  const networks = data?.getNetworks.networks ?? [];
  const existing = networks.find((n) => n.isDefault) ?? networks[0];
  const guarded = !isAddFlow && Boolean(existing);

  useEffect(() => {
    if (!isAddFlow && existing) {
      // Keep the dashboard's selected network in sync (replaces the seed
      // default) so screens never query a non-existent network id.
      setNetworkId(existing.id);
      router.replace(
        stepUrl('/configure/install', { flow, networkid: existing.id }),
      );
    }
  }, [isAddFlow, existing, flow, router, setNetworkId]);

  const [addNetwork, { loading: saving }] = useAddNetworkMutation({
    refetchQueries: [
      { query: OnboardingStatusDocument },
      { query: GetNetworksDocument },
    ],
    onCompleted: (res) => {
      setNetworkId(res.addNetwork.id);
      router.push(
        stepUrl('/configure/install', { flow, networkid: res.addNetwork.id }),
      );
    },
    onError: (err) => setSubmitError(err.message),
  });

  const onSubmit = (values: FormValues) => {
    setSubmitError(null);
    void addNetwork({
      // First network becomes the default; additional ones don't steal it.
      variables: {
        data: { name: values.name, isDefault: networks.length === 0 },
      },
    });
  };

  if (loading || guarded) {
    return (
      <>
        <Skeleton width="60%" height={42} />
        <Skeleton width="100%" height={28} />
        <Skeleton width="100%" height={56} sx={{ mt: 2 }} />
      </>
    );
  }

  return (
    <>
      <h1 className="cfg-title">Name your network</h1>
      <p className="cfg-copy">
        A network is made up of one or more sites of Ukama hardware, allowing
        you to connect to the cellular internet.
      </p>
      <form
        onSubmit={(e) => void handleSubmit(onSubmit)(e)}
        className="cfg-fields"
        style={{ display: 'flex', flexDirection: 'column', flex: 1 }}
      >
        <Field label="Network name" required error={errors.name?.message}>
          <TextInput
            placeholder="network-name"
            invalid={Boolean(errors.name)}
            autoFocus
            {...register('name')}
          />
        </Field>
        {submitError && <p className="cfg-error">{submitError}</p>}
        <ConfigureActions
          nextLabel="Create network"
          onNext={() => void handleSubmit(onSubmit)()}
          busy={saving}
        />
      </form>
    </>
  );
}
