/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Graphs_Type } from '@/client/graphql/generated/subscriptions';
import PubSub from 'pubsub-js';

interface IMetricSubscription {
  url: string;
  key: string;
  from: number;
  nodeId: string;
  userId: string;
  orgName: string;
  type: Graphs_Type;
}

interface SSEEvent {
  data?: string;
  id?: string;
  event?: string;
}

const activeSubscriptions = new Map<string, AbortController>();

const RETRY_DELAY_MS = 2000;
const ACTIVITY_TIMEOUT_MS = 5 * 60 * 1000;
const MAX_BUFFER_SIZE = 25000;

function parseEvent(eventStr: string): SSEEvent | null {
  if (!eventStr || eventStr.startsWith(':')) return null;

  const event: SSEEvent = {};
  const lines = eventStr.split('\n');

  for (const line of lines) {
    if (line.startsWith('data:')) {
      event.data = line.slice(5).trim();
    } else if (line.startsWith('id:')) {
      event.id = line.slice(3).trim();
    } else if (line.startsWith('event:')) {
      event.event = line.slice(6).trim();
    }
  }

  return event.data || event.id || event.event ? event : null;
}

export default function MetricSubscription({
  url,
  key,
  from,
  type,
  userId,
  nodeId,
  orgName,
}: IMetricSubscription) {
  if (activeSubscriptions.has(key)) {
    activeSubscriptions.get(key)!.abort();
    activeSubscriptions.delete(key);
  }

  const queryParams = new URLSearchParams({
    query:
      'subscription GetMetricByTabSub($data:SubMetricByTabInput!){getMetricByTabSub(data:$data){msg nodeId success type value}}',
    variables: JSON.stringify({
      data: { nodeId, orgName, type, userId, from },
    }),
    operationName: 'GetMetricByTabSub',
    extensions: '{}',
  });

  const fullUrl = `${url}/graphql?${queryParams.toString()}`;

  const controller = new AbortController();
  activeSubscriptions.set(key, controller);

  const myHeaders = new Headers({
    'Cache-Control': 'no-cache',
    Connection: 'keep-alive',
    Pragma: 'no-cache',
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'same-origin',
    accept: 'text/event-stream',
  });

  const handleBeforeUnload = () => {
    if (activeSubscriptions.has(key)) {
      activeSubscriptions.get(key)!.abort();
      activeSubscriptions.delete(key);
    }
  };
  window.addEventListener('beforeunload', handleBeforeUnload);

  async function startSSE() {
    try {
      const response = await fetch(fullUrl, {
        method: 'GET',
        headers: myHeaders,
        signal: controller.signal,
      });

      if (!response.ok || !response.body) {
        throw new Error(`HTTP Error ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder('utf-8');
      let buffer = '';

      let activityTimeout: ReturnType<typeof setTimeout> | null = null;
      const resetActivityTimeout = () => {
        if (activityTimeout) clearTimeout(activityTimeout);
        activityTimeout = setTimeout(() => {
          controller.abort();
          activeSubscriptions.delete(key);
        }, ACTIVITY_TIMEOUT_MS);
      };
      resetActivityTimeout();

      while (true) {
        if (controller.signal.aborted) break;

        const { value, done } = await reader.read();
        if (done) break;

        resetActivityTimeout();

        const chunk = decoder.decode(value, { stream: true });
        if (buffer.length + chunk.length > MAX_BUFFER_SIZE) {
          buffer = buffer.slice(-MAX_BUFFER_SIZE / 2);
        }
        buffer += chunk;

        let eventBoundary = buffer.indexOf('\n\n');
        while (eventBoundary !== -1) {
          const eventStr = buffer.slice(0, eventBoundary).trim();
          buffer = buffer.slice(eventBoundary + 2);

          const parsedEvent = parseEvent(eventStr);
          if (parsedEvent?.data) {
            PubSub.publish(key, parsedEvent.data);
          }

          eventBoundary = buffer.indexOf('\n\n');
        }
      }
    } catch (error) {
      if (error instanceof Error && error.name === 'AbortError') return;
      if (activeSubscriptions.has(key)) {
        setTimeout(startSSE, RETRY_DELAY_MS);
      }
    } finally {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    }
  }

  startSSE();

  return {
    cancel: () => {
      if (activeSubscriptions.has(key)) {
        activeSubscriptions.get(key)!.abort();
        activeSubscriptions.delete(key);
        window.removeEventListener('beforeunload', handleBeforeUnload);
      }
    },
  };
}
