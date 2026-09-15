/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Runs site creation as three visible steps.
 *
 *   1. checks the site's three nodes are online and READY;
 *   2. calls addSite, which configures those nodes and only writes the site
 *      once all three report OPERATIONAL (registry/site provisioning barrier,
 *      up to ~3 minutes);
 *   3. waits for the site to appear in the network's site list.
 *
 * Step 3 polls rather than trusting the mutation's own result: the barrier
 * outlives the request that started it (the BFF aborts upstream calls at
 * HTTP_TIMEOUT_MS, and ingress cuts earlier than the barrier's ceiling), so a
 * transport failure there means "still provisioning", not "failed". A
 * GraphQL error is the backend's own verdict and stops the run.
 */
import { useCallback, useRef, useState } from 'react';
import { useApolloClient } from '@apollo/client';

import {
  GetNodesDocument,
  type GetNodesQuery,
  type GetNodesQueryVariables,
} from '@/client/graphql/nodes.generated';
import {
  GetSitesDocument,
  type GetSitesQuery,
  type GetSitesQueryVariables,
} from '@/client/graphql/sites.generated';
import type { AddSiteInputDto } from '@/client/graphql/types';
import { useAddSiteMutation } from '@/client/graphql/sites.generated';
import { isOnlineAndReady, siteUnits, stateLabel } from './detectSites';

export type StepKey = 'nodes' | 'create' | 'confirm';
export type StepState = 'pending' | 'active' | 'done' | 'failed';

export interface CreateStep {
  key: StepKey;
  title: string;
  /** Live detail under the title: which node is lagging, elapsed time, error. */
  detail: string;
  state: StepState;
}

/** Node readiness poll: nodes settle in seconds, so give up reasonably fast. */
const NODES_TIMEOUT_MS = 60_000;
const NODES_POLL_MS = 3_000;
/** Provisioning ceiling is 3 attempts x 60s, plus room for the site to land. */
const CONFIRM_TIMEOUT_MS = 240_000;
const CONFIRM_POLL_MS = 5_000;

const STEP_TITLES: Record<StepKey, string> = {
  nodes: 'Checking your nodes',
  create: 'Creating your site',
  confirm: 'Confirming your site',
};

const initialSteps = (): CreateStep[] =>
  (['nodes', 'create', 'confirm'] as StepKey[]).map((key) => ({
    key,
    title: STEP_TITLES[key],
    detail: '',
    state: 'pending',
  }));

const wait = (ms: number) => new Promise((r) => setTimeout(r, ms));

export interface CreateSiteRun {
  steps: CreateStep[];
  running: boolean;
  /** Set when a step failed; the steps carry the detail. */
  error: string | null;
  run: (input: AddSiteInputDto, towerId: string) => Promise<boolean>;
  reset: () => void;
}

export function useCreateSite(networkId: string): CreateSiteRun {
  const client = useApolloClient();
  const [addSite] = useAddSiteMutation();
  const [steps, setSteps] = useState<CreateStep[]>(initialSteps);
  const [running, setRunning] = useState(false);
  const [error, setError] = useState<string | null>(null);
  // Guards a second run from the same click-through.
  const active = useRef(false);

  const patch = useCallback((key: StepKey, next: Partial<CreateStep>) => {
    setSteps((prev) =>
      prev.map((s) => (s.key === key ? { ...s, ...next } : s)),
    );
  }, []);

  const reset = useCallback(() => {
    setSteps(initialSteps());
    setError(null);
    setRunning(false);
    active.current = false;
  }, []);

  /** Step 1: all three units online and READY. */
  const checkNodes = useCallback(
    async (towerId: string): Promise<boolean> => {
      patch('nodes', { state: 'active', detail: 'Reading node states…' });
      const deadline = Date.now() + NODES_TIMEOUT_MS;

      for (;;) {
        const { data } = await client.query<
          GetNodesQuery,
          GetNodesQueryVariables
        >({
          query: GetNodesDocument,
          variables: { data: {} },
          fetchPolicy: 'network-only',
          errorPolicy: 'all',
        });

        const units = siteUnits(data?.getNodes.nodes ?? [], towerId);
        const missing = units.missing;
        const notReady = units.found.filter((n) => !isOnlineAndReady(n));

        if (!missing.length && !notReady.length && units.found.length === 3) {
          patch('nodes', {
            state: 'done',
            detail: 'Tower, amplifier and controller are online and ready.',
          });
          return true;
        }

        if (Date.now() >= deadline) {
          const why = missing.length
            ? `${missing.join(', ')} not found for this site.`
            : notReady
                .map((n) => `${n.type.toLowerCase()} is ${stateLabel(n)}`)
                .join(', ');
          patch('nodes', { state: 'failed', detail: why });
          setError(`Your nodes are not ready yet. ${why}`);
          return false;
        }

        patch('nodes', {
          state: 'active',
          detail: missing.length
            ? `Waiting for ${missing.join(', ')}…`
            : `Waiting for ${notReady
                .map((n) => `${n.type.toLowerCase()} (${stateLabel(n)})`)
                .join(', ')}…`,
        });
        await wait(NODES_POLL_MS);
      }
    },
    [client, patch],
  );

  /** Step 3: the site row lands only after all three nodes are operational. */
  const confirmSite = useCallback(
    async (name: string): Promise<boolean> => {
      const started = Date.now();
      const deadline = started + CONFIRM_TIMEOUT_MS;

      for (;;) {
        const { data } = await client.query<
          GetSitesQuery,
          GetSitesQueryVariables
        >({
          query: GetSitesDocument,
          variables: { data: { networkId } },
          fetchPolicy: 'network-only',
          errorPolicy: 'all',
        });

        const found = (data?.getSites.sites ?? []).some(
          (s) => s.name.toLowerCase() === name.toLowerCase(),
        );
        if (found) {
          patch('confirm', { state: 'done', detail: `“${name}” is ready.` });
          return true;
        }

        if (Date.now() >= deadline) {
          patch('confirm', {
            state: 'failed',
            detail: 'The site has not appeared yet.',
          });
          setError(
            'Your nodes are still being configured. The site will appear once they finish, or you can try again.',
          );
          return false;
        }

        const elapsed = Math.round((Date.now() - started) / 1000);
        patch('confirm', {
          state: 'active',
          detail: `Configuring your nodes… ${elapsed}s`,
        });
        await wait(CONFIRM_POLL_MS);
      }
    },
    [client, networkId, patch],
  );

  const run = useCallback(
    async (input: AddSiteInputDto, towerId: string): Promise<boolean> => {
      if (active.current) return false;
      active.current = true;
      setSteps(initialSteps());
      setError(null);
      setRunning(true);

      try {
        if (!(await checkNodes(towerId))) return false;

        patch('create', {
          state: 'active',
          detail: 'Sending your site configuration…',
        });
        try {
          await addSite({ variables: { data: input } });
        } catch (e) {
          const err = e as {
            graphQLErrors?: readonly { message: string }[];
            message?: string;
          };
          // The backend rejected it: its verdict is final, so stop here.
          if (err.graphQLErrors?.length) {
            const msg = err.graphQLErrors[0]!.message;
            patch('create', { state: 'failed', detail: msg });
            setError(msg);
            return false;
          }
          // Transport gave up while the barrier runs on. Confirm by polling.
          patch('create', {
            state: 'done',
            detail: 'Sent. Your nodes are being configured.',
          });
          return await confirmSite(input.name);
        }

        patch('create', { state: 'done', detail: 'Site configuration sent.' });
        return await confirmSite(input.name);
      } finally {
        setRunning(false);
        active.current = false;
      }
    },
    [addSite, checkNodes, confirmSite, patch],
  );

  return { steps, running, error, run, reset };
}
