/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Site step 2 of 2 — "Configure site settings". Confirms the switch, power,
 * and backhaul components for the named site, then creates it. Self-guards:
 * the previous step must have provided a site name + tower (nid); otherwise
 * we send the user back to /configure/site to (re)detect and name it.
 */
'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { useGetComponentsByUserIdQuery } from '@/client/graphql/components.generated';
import { useGetNodesQuery } from '@/client/graphql/nodes.generated';
import { OnboardingStatusDocument } from '@/client/graphql/onboarding-status.generated';
import { useAddSiteMutation } from '@/client/graphql/sites.generated';
import { Component_Type } from '@/client/graphql/types';
import { Field, SelectInput } from '@/components/form/FormField';
import ConfigureActions from '../../_components/ConfigureActions';
import { parseCoords } from '../../_components/coords';
import { stepUrl, useConfigureParams } from '../../_components/state';

const schema = z.object({
  powerId: z.string().min(1, 'Select a power component'),
  backhaulId: z.string().min(1, 'Select a backhaul component'),
  switchId: z.string().min(1, 'Select a switch component'),
});

type FormValues = z.infer<typeof schema>;

function StepSkeleton() {
  return (
    <>
      <Skeleton width="60%" height={42} />
      <Skeleton width="100%" height={28} />
      <Skeleton width="100%" height={56} sx={{ mt: 2 }} />
      <Skeleton width="100%" height={56} />
    </>
  );
}

export default function ConfigureSiteSettingsPage() {
  const router = useRouter();
  const params = useConfigureParams();
  const [submitError, setSubmitError] = useState<string | null>(null);

  // Self-guard: this step needs a named site + tower from the previous step.
  const incomplete = !params.nid || !params.sitename || !params.networkid;
  useEffect(() => {
    if (incomplete) {
      router.replace(stepUrl('/configure/site', { flow: params.flow }));
    }
  }, [incomplete, params.flow, router]);

  // Re-fetch the tower for its coordinates (resumable, no cross-step state).
  const { data: nodesData, loading: nodesLoading } = useGetNodesQuery({
    variables: { data: {} },
    fetchPolicy: 'cache-first',
    skip: incomplete,
  });
  const tower = nodesData?.getNodes.nodes.find((n) => n.id === params.nid);

  const { data: compData, loading: compLoading } =
    useGetComponentsByUserIdQuery({
      variables: { data: { category: Component_Type.All } },
      skip: incomplete,
    });
  const components = useMemo(
    () => compData?.getComponentsByUserId.components ?? [],
    [compData],
  );
  const byCategory = (cat: Component_Type) =>
    components.filter((c) => c.category === cat);
  const power = byCategory(Component_Type.Power);
  const backhaul = byCategory(Component_Type.Backhaul);
  const switches = byCategory(Component_Type.Switch);
  const spectrum = byCategory(Component_Type.Spectrum);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { powerId: '', backhaulId: '', switchId: '' },
  });

  // Auto-select unambiguous choices (legacy parity).
  useEffect(() => {
    if (power.length === 1 && power[0]) setValue('powerId', power[0].id);
    if (backhaul.length === 1 && backhaul[0])
      setValue('backhaulId', backhaul[0].id);
    if (switches.length === 1 && switches[0])
      setValue('switchId', switches[0].id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [components.length]);

  const [addSite, { loading: saving }] = useAddSiteMutation({
    refetchQueries: [{ query: OnboardingStatusDocument }],
    onCompleted: () => router.push(stepUrl('/configure/sims', { flow: params.flow })),
    onError: (err) => setSubmitError(err.message),
  });

  const onSubmit = (values: FormValues) => {
    setSubmitError(null);
    // access = inventory component whose partNumber is the tower id (legacy).
    const access = components.find(
      (c) =>
        c.category === Component_Type.Access && c.partNumber === params.nid,
    );
    if (!access) {
      setSubmitError(
        'Your node was not found in your inventory (no access component). Please contact support.',
      );
      return;
    }
    const spectrumId = spectrum[0]?.id;
    if (!spectrumId) {
      setSubmitError(
        'No spectrum component found in your inventory. Please contact support.',
      );
      return;
    }
    // Normalize possibly-swapped coordinates before persisting (same rule the
    // name step uses for the map/address).
    const coords = parseCoords(tower?.latitude, tower?.longitude);
    void addSite({
      variables: {
        data: {
          name: params.sitename,
          network_id: params.networkid,
          access_id: access.id,
          power_id: values.powerId,
          backhaul_id: values.backhaulId,
          switch_id: values.switchId,
          spectrum_id: spectrumId,
          location: params.location,
          latitude: String(coords?.lat ?? tower?.latitude ?? '0'),
          longitude: String(coords?.lng ?? tower?.longitude ?? '0'),
          install_date: new Date().toISOString(),
        },
      },
    });
  };

  if (incomplete || nodesLoading || compLoading) return <StepSkeleton />;

  const toOptions = (list: typeof power) =>
    list.map((c) => ({ value: c.id, label: c.description || c.partNumber }));

  return (
    <>
      <h1 className="cfg-title">Configure site settings</h1>
      <p className="cfg-copy">
        If you used the default Ukama hardware, these are already selected —
        just continue. If you used your own switch, power, or backhaul, pick
        the matching option. Note: we can&apos;t track real-time KPIs for
        custom components.
      </p>
      <form
        onSubmit={(e) => void handleSubmit(onSubmit)(e)}
        className="cfg-fields"
        style={{ display: 'flex', flexDirection: 'column', flex: 1 }}
      >
        <Field label="Node">
          <div className="cfg-readonly">{params.nid}</div>
        </Field>
        <Field label="Switch" required error={errors.switchId?.message}>
          <SelectInput
            placeholder="Select switch component"
            options={toOptions(switches)}
            invalid={Boolean(errors.switchId)}
            {...register('switchId')}
          />
        </Field>
        <Field label="Backhaul" required error={errors.backhaulId?.message}>
          <SelectInput
            placeholder="Select backhaul component"
            options={toOptions(backhaul)}
            invalid={Boolean(errors.backhaulId)}
            {...register('backhaulId')}
          />
        </Field>
        <Field label="Power" required error={errors.powerId?.message}>
          <SelectInput
            placeholder="Select power component"
            options={toOptions(power)}
            invalid={Boolean(errors.powerId)}
            {...register('powerId')}
          />
        </Field>
        {submitError && <p className="cfg-error">{submitError}</p>}
        <ConfigureActions
          nextLabel="Create site"
          onNext={() => void handleSubmit(onSubmit)()}
          busy={saving}
        />
      </form>
    </>
  );
}
