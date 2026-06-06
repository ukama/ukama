/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Site setup: detect an installed node, name the site, confirm the
 * inventory components (auto-selected when unambiguous), create the site.
 * Self-guards: no network → back to the network step; no node detected →
 * friendly retry/skip state (the physical install may still be in flight).
 */
'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { useGetComponentsByUserIdQuery } from '@/client/graphql/components.generated';
import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { useGetNodesQuery } from '@/client/graphql/nodes.generated';
import { OnboardingStatusDocument } from '@/client/graphql/onboarding-status.generated';
import { useAddSiteMutation } from '@/client/graphql/sites.generated';
import { Component_Type } from '@/client/graphql/types';
import { Field, SelectInput, TextInput } from '@/components/form/FormField';
import ConfigureActions from '../_components/ConfigureActions';
import SiteReadinessChecklist from '../_components/SiteReadinessChecklist';
import {
  computeSiteReadiness,
  detectReadySites,
} from '../_components/detectSites';
import { stepUrl, useConfigureParams } from '../_components/state';

/** Poll cadence for node detection while the site is coming online. */
const NODE_POLL_MS = 15_000;

const schema = z.object({
  name: z
    .string()
    .min(3, 'At least 3 characters')
    .max(40, 'At most 40 characters')
    .regex(
      /^[a-z0-9-]+$/,
      'Lowercase letters, numbers, and "-" only (no spaces).',
    ),
  location: z.string().min(3, 'Location is required'),
  nodeId: z.string().min(1, 'Select the installed node'),
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

export default function ConfigureSitePage() {
  const router = useRouter();
  const params = useConfigureParams();
  const [submitError, setSubmitError] = useState<string | null>(null);

  // --- network guard: resolve a networkid or send back to the network step.
  const { data: netData, loading: netLoading } = useGetNetworksQuery({
    skip: Boolean(params.networkid),
  });
  const fallbackNet =
    netData?.getNetworks.networks.find((n) => n.isDefault) ??
    netData?.getNetworks.networks[0];
  const networkId = params.networkid || (fallbackNet?.id ?? '');

  useEffect(() => {
    if (!params.networkid && !netLoading && netData && !fallbackNet) {
      router.replace(stepUrl('/configure/network', { flow: params.flow }));
    }
  }, [params.networkid, params.flow, netLoading, netData, fallbackNet, router]);

  // --- node detection: a site needs a full, powered, unconfigured trio
  // (tower + amplifier + controller). Poll getNodes every 15s while the site
  // is coming online; stop once it's ready.
  const [lastFetched, setLastFetched] = useState<Date | null>(null);
  const {
    data: nodesData,
    loading: nodesLoading,
    refetch: refetchNodes,
    startPolling,
    stopPolling,
  } = useGetNodesQuery({
    variables: { data: {} },
    fetchPolicy: 'network-only',
    notifyOnNetworkStatusChange: true,
    // Stamp the time each poll/fetch settles (fires on every poll cycle).
    onCompleted: () => setLastFetched(new Date()),
    onError: () => setLastFetched(new Date()),
  });
  const nodes = useMemo(() => nodesData?.getNodes.nodes ?? [], [nodesData]);
  const readiness = useMemo(
    () => computeSiteReadiness(nodes, params.nid || undefined),
    [nodes, params.nid],
  );
  const readySites = useMemo(() => {
    const sites = detectReadySites(nodes);
    // A deep-linked node id pins the flow to that specific tower's site.
    return params.nid ? sites.filter((s) => s.tower.id === params.nid) : sites;
  }, [nodes, params.nid]);

  // Poll while the site isn't ready; stop polling once it is.
  useEffect(() => {
    if (readiness.ready) stopPolling();
    else startPolling(NODE_POLL_MS);
    return () => stopPolling();
  }, [readiness.ready, startPolling, stopPolling]);

  // --- inventory components.
  const { data: compData, loading: compLoading } =
    useGetComponentsByUserIdQuery({
      variables: { data: { category: Component_Type.All } },
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
    defaultValues: {
      name: '',
      location: '',
      nodeId: '',
      powerId: '',
      backhaulId: '',
      switchId: '',
    },
  });

  // Auto-select unambiguous choices (legacy parity). Node = detected tower.
  useEffect(() => {
    const first = readySites[0];
    if (first) setValue('nodeId', first.tower.id);
    if (power.length === 1 && power[0]) setValue('powerId', power[0].id);
    if (backhaul.length === 1 && backhaul[0])
      setValue('backhaulId', backhaul[0].id);
    if (switches.length === 1 && switches[0])
      setValue('switchId', switches[0].id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [readySites.length, components.length]);

  const [addSite, { loading: saving }] = useAddSiteMutation({
    refetchQueries: [{ query: OnboardingStatusDocument }],
    onCompleted: () => {
      router.push(stepUrl('/configure/sims', { flow: params.flow })); // networkid no longer needed
    },
    onError: (err) => setSubmitError(err.message),
  });

  const onSubmit = (values: FormValues) => {
    setSubmitError(null);
    const node = readySites.find((s) => s.tower.id === values.nodeId)?.tower;
    // access = inventory component whose partNumber is the node id (legacy).
    const access = components.find(
      (c) =>
        c.category === Component_Type.Access && c.partNumber === values.nodeId,
    );
    if (!access) {
      setSubmitError(
        'The selected node was not found in your inventory (no access component). Please contact support.',
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
    void addSite({
      variables: {
        data: {
          name: values.name,
          network_id: networkId,
          access_id: access.id,
          power_id: values.powerId,
          backhaul_id: values.backhaulId,
          switch_id: values.switchId,
          spectrum_id: spectrumId,
          location: values.location,
          latitude: String(node?.latitude ?? '0'),
          longitude: String(node?.longitude ?? '0'),
          install_date: new Date().toISOString(),
        },
      },
    });
  };

  // Initial load only: once polling starts we keep the checklist on screen.
  if (compLoading || (nodesLoading && nodes.length === 0)) {
    return <StepSkeleton />;
  }

  if (!readiness.ready) {
    return (
      <>
        <h1 className="cfg-title">Bring your site online</h1>
        <p className="cfg-copy">
          A site is made up of three Ukama units working together. Follow the
          steps below — we&apos;re checking automatically every 15 seconds and
          will continue on our own once your site is online. You can also skip
          and finish later.
        </p>
        <SiteReadinessChecklist
          readiness={readiness}
          lastFetched={lastFetched}
          refreshing={nodesLoading}
          onRefresh={() => void refetchNodes()}
        />
        <div className="cfg-actions">
          <Button variant="text" sx={{ color: 'var(--uk-ink-2)' }} onClick={() => router.push('/')}>
            Skip for now
          </Button>
        </div>
      </>
    );
  }

  const nodeOptions = readySites.map((s) => ({
    value: s.tower.id,
    label: s.tower.name ? `${s.tower.name} (${s.tower.id})` : s.tower.id,
  }));
  const toOptions = (list: typeof power) =>
    list.map((c) => ({ value: c.id, label: c.description || c.partNumber }));

  return (
    <>
      <h1 className="cfg-title">Name your site</h1>
      <p className="cfg-copy">
        We found your installed hardware. Confirm the details below to create
        your first site.
      </p>
      <form
        onSubmit={(e) => void handleSubmit(onSubmit)(e)}
        className="cfg-fields"
        style={{ display: 'flex', flexDirection: 'column', flex: 1 }}
      >
        <Field label="Site name" required error={errors.name?.message}>
          <TextInput
            placeholder="site-name"
            invalid={Boolean(errors.name)}
            autoFocus
            {...register('name')}
          />
        </Field>
        <Field label="Location" required error={errors.location?.message}>
          <TextInput
            placeholder="Street, city, country"
            invalid={Boolean(errors.location)}
            {...register('location')}
          />
        </Field>
        <Field label="Node" required error={errors.nodeId?.message}>
          <SelectInput
            options={nodeOptions}
            invalid={Boolean(errors.nodeId)}
            {...register('nodeId')}
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
        <Field label="Backhaul" required error={errors.backhaulId?.message}>
          <SelectInput
            placeholder="Select backhaul component"
            options={toOptions(backhaul)}
            invalid={Boolean(errors.backhaulId)}
            {...register('backhaulId')}
          />
        </Field>
        <Field label="Switch" required error={errors.switchId?.message}>
          <SelectInput
            placeholder="Select switch component"
            options={toOptions(switches)}
            invalid={Boolean(errors.switchId)}
            {...register('switchId')}
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
