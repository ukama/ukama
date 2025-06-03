/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Stats_Type } from '@/client/graphql/generated/subscriptions';
import PubSub from 'pubsub-js';

interface IMetricStatBySiteSubscription {
  url: string;
  key: string;
  from: number;
  siteIds: string[];
  userId: string;
  orgName: string;
  type: Stats_Type;
  nodeIds?: string[];
}

const activeSubscriptions = new Map<string, AbortController>();
function parseEvent(eventStr: string) {
  if (!eventStr || eventStr.startsWith(':')) return null;

  const event: any = {};
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
export default function MetricStatBySiteSubscription(
  params: IMetricStatBySiteSubscription,
) {
  const { url, key, from, type, userId, siteIds, nodeIds, orgName } = params;

  if (activeSubscriptions.has(key)) {
    const existingController = activeSubscriptions.get(key)!;
    existingController.abort();
    activeSubscriptions.delete(key);
  }

  const queryParams = new URLSearchParams();
  const variables = {
    data: {
      siteIds,
      orgName,
      type,
      userId,
      from,
      ...(nodeIds && nodeIds.length > 0 && { nodeIds }),
    },
  };

  queryParams.set(
    'query',
    'subscription SiteMetricStatSub($data: SubSiteMetricsStatInput!) { getSiteMetricStatSub(data: $data) { msg siteId nodeId success type value } }',
  );
  queryParams.set('variables', JSON.stringify(variables));
  queryParams.set('operationName', 'SiteMetricStatSub');

  const fullUrl = `${url}/graphql?${queryParams.toString()}`;

  const controller = new AbortController();
  activeSubscriptions.set(key, controller);

  const myHeaders = new Headers();
  myHeaders.append('Cache-Control', 'no-cache');
  myHeaders.append('Connection', 'keep-alive');
  myHeaders.append('Pragma', 'no-cache');
  myHeaders.append('Sec-Fetch-Dest', 'empty');
  myHeaders.append('Sec-Fetch-Mode', 'cors');
  myHeaders.append('Sec-Fetch-Site', 'same-origin');
  myHeaders.append('accept', 'text/event-stream');

  const requestOptions = {
    method: 'GET',
    headers: myHeaders,
    signal: controller.signal,
  };

  const handleBeforeUnload = () => {
    if (activeSubscriptions.has(key)) {
      const controller = activeSubscriptions.get(key)!;
      controller.abort();
      activeSubscriptions.delete(key);
    }
  };
  window.addEventListener('beforeunload', handleBeforeUnload);

  const MAX_BUFFER_SIZE = 25000;

  async function startSSE() {
    try {
      const response = await fetch(fullUrl, requestOptions);

      if (!response.ok || !response.body) {
        throw new Error(`HTTP Error ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder('utf-8');
      let buffer = '';

      let activityTimeout: any = null;
      const resetActivityTimeout = () => {
        if (activityTimeout) clearTimeout(activityTimeout);
        activityTimeout = setTimeout(
          () => {
            controller.abort();
            activeSubscriptions.delete(key);
          },
          5 * 60 * 1000,
        );
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
    } catch (error: any) {
      if (error.name === 'AbortError') {
      } else {
        if (activeSubscriptions.has(key)) {
          setTimeout(() => {
            startSSE();
          }, 2000);
        }
      }
    } finally {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    }
  }

  startSSE();

  return {
    cancel: () => {
      if (activeSubscriptions.has(key)) {
        const controller = activeSubscriptions.get(key)!;
        controller.abort();
        activeSubscriptions.delete(key);
        window.removeEventListener('beforeunload', handleBeforeUnload);
      }
    },
  };
}
