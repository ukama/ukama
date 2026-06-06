/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Site step 1 of 2 — "Name your site". Detects the installed node trio (a
 * guided, polled checklist until the site is online), then shows the tower's
 * location on a map with its address and asks for a site name. The component
 * selection + creation happens on the next step (/configure/site/settings).
 * Self-guard: no network → back to the network step.
 */
'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import Button from '@mui/material/Button';
import Skeleton from '@mui/material/Skeleton';
import { useRouter } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { useGetNetworksQuery } from '@/client/graphql/networks.generated';
import { useGetNodesQuery } from '@/client/graphql/nodes.generated';
import { Field, TextInput } from '@/components/form/FormField';
import ConfigureActions from '../_components/ConfigureActions';
import SiteLocationMap from '../_components/SiteLocationMap';
import SiteReadinessChecklist from '../_components/SiteReadinessChecklist';
import { parseCoords } from '../_components/coords';
import { computeSiteReadiness } from '../_components/detectSites';
import { stepUrl, useConfigureParams } from '../_components/state';
import { useReverseGeocode } from '../_components/useReverseGeocode';

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

export default function ConfigureSiteNamePage() {
  const router = useRouter();
  const params = useConfigureParams();

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

  // --- node detection: poll getNodes every 15s while coming online.
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
    onCompleted: () => setLastFetched(new Date()),
    onError: () => setLastFetched(new Date()),
  });
  const nodes = useMemo(() => nodesData?.getNodes.nodes ?? [], [nodesData]);
  const readiness = useMemo(
    () => computeSiteReadiness(nodes, params.nid || undefined),
    [nodes, params.nid],
  );
  const tower = readiness.towerNode;

  useEffect(() => {
    if (readiness.ready) stopPolling();
    else startPolling(NODE_POLL_MS);
    return () => stopPolling();
  }, [readiness.ready, startPolling, stopPolling]);

  // --- tower location → address (coords may be stored swapped; normalize).
  const coords = tower
    ? parseCoords(tower.latitude, tower.longitude)
    : null;
  const { address } = useReverseGeocode(
    coords?.lat ?? null,
    coords?.lng ?? null,
  );

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: '' },
  });

  const onSubmit = (values: FormValues) => {
    if (!tower) return;
    router.push(
      stepUrl('/configure/site/settings', {
        flow: params.flow,
        networkid: networkId,
        nid: tower.id,
        sitename: values.name,
        location: address,
      }),
    );
  };

  // Initial load only: once polling starts we keep the checklist on screen.
  if (nodesLoading && nodes.length === 0) return <StepSkeleton />;

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
          <Button
            variant="text"
            sx={{ color: 'var(--uk-ink-2)' }}
            onClick={() => router.push('/')}
          >
            Skip for now
          </Button>
        </div>
      </>
    );
  }

  return (
    <>
      <h1 className="cfg-title">Name your site</h1>
      <p className="cfg-copy">
        We found your site. Confirm the location below and give it a name for
        easy reference.
      </p>
      <form
        onSubmit={(e) => void handleSubmit(onSubmit)(e)}
        className="cfg-fields"
        style={{ display: 'flex', flexDirection: 'column', flex: 1 }}
      >
        {coords && <SiteLocationMap lat={coords.lat} lng={coords.lng} />}
        <Field label="Site location">
          <div className="cfg-readonly">{address || 'Resolving location…'}</div>
        </Field>
        <Field label="Site name" required error={errors.name?.message}>
          <TextInput
            placeholder="site-name"
            invalid={Boolean(errors.name)}
            autoFocus
            {...register('name')}
          />
        </Field>
        <ConfigureActions
          nextLabel="Name site"
          onNext={() => void handleSubmit(onSubmit)()}
        />
      </form>
    </>
  );
}
