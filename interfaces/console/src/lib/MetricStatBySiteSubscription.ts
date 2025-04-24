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

interface SubscriptionController {
  controller: AbortController;
  lastActivity: number;
}

const activeSubscriptions = new Map<string, SubscriptionController>();

function parseEvent(eventStr: string): {
  data?: string;
  id?: string;
  event?: string;
} {
  const event: { data?: string; id?: string; event?: string } = {};
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

  return event;
}

export default function MetricStatBySiteSubscription(
  params: IMetricStatBySiteSubscription,
): { cancel: () => void } {
  const { url, key, from, type, userId, siteIds, nodeIds, orgName } = params;

  if (activeSubscriptions.has(key)) {
    const { controller } = activeSubscriptions.get(key)!;
    controller.abort();
    activeSubscriptions.delete(key);
  }

  const controller = new AbortController();

  activeSubscriptions.set(key, {
    controller,
    lastActivity: Date.now(),
  });

  const headers = new Headers({
    'Cache-Control': 'no-cache',
    Connection: 'keep-alive',
    Pragma: 'no-cache',
    accept: 'text/event-stream',
  });

  const siteIdsParam = encodeURIComponent(JSON.stringify(siteIds));
  let fullUrl = `${url}/graphql?query=subscription+SiteMetricStatSub%28%24data%3ASubSiteMetricsStatInput%21%29%7BgetSiteMetricStatSub%28data%3A%24data%29%7Bmsg+siteId+nodeId+success+type+value%7D%7D&variables=%7B%22data%22%3A%7B%22siteIds%22%3A${siteIdsParam}%2C%22orgName%22%3A%22${orgName}%22%2C%22type%22%3A%22${type}%22%2C%22userId%22%3A%22${userId}%22%2C%22from%22%3A${from}`;

  if (nodeIds && nodeIds.length > 0) {
    const nodeIdsParam = encodeURIComponent(JSON.stringify(nodeIds));
    fullUrl += `%2C%22nodeIds%22%3A${nodeIdsParam}`;
  }

  fullUrl += `%7D%7D&operationName=SiteMetricStatSub&extensions=%7B%7D`;

  const unloadListener = () => {
    if (activeSubscriptions.has(key)) {
      activeSubscriptions.get(key)!.controller.abort();
      activeSubscriptions.delete(key);
    }
  };

  window.addEventListener('beforeunload', unloadListener);

  const cleanupToken = PubSub.subscribe(`cleanup-${key}`, () => {
    if (activeSubscriptions.has(key)) {
      activeSubscriptions.get(key)!.controller.abort();
      activeSubscriptions.delete(key);
      window.removeEventListener('beforeunload', unloadListener);
      PubSub.unsubscribe(cleanupToken);
    }
  });

  (async () => {
    try {
      const response = await fetch(fullUrl, {
        method: 'GET',
        headers,
        signal: controller.signal,
      });

      if (!response || !response.ok) {
        console.error(
          `Network error for subscription ${key}: ${response?.status}`,
        );
        return;
      }

      const reader = response.body?.getReader();
      if (!reader) {
        console.error('Could not get reader from response body');
        return;
      }

      const decoder = new TextDecoder('utf-8');
      let buffer = '';
      const MAX_BUFFER_SIZE = 50000;

      while (!controller.signal.aborted) {
        const { value, done } = await reader.read();

        if (done) {
          break;
        }

        if (activeSubscriptions.has(key)) {
          activeSubscriptions.get(key)!.lastActivity = Date.now();
        }

        buffer += decoder.decode(value, { stream: true });

        if (buffer.length > MAX_BUFFER_SIZE) {
          buffer = buffer.slice(-MAX_BUFFER_SIZE / 2);
        }

        let eventBoundary = buffer.indexOf('\n\n');
        while (eventBoundary !== -1) {
          const eventStr = buffer.slice(0, eventBoundary).trim();
          buffer = buffer.slice(eventBoundary + 2);

          if (eventStr !== ':') {
            const event = parseEvent(eventStr);
            if (event.data) {
              PubSub.publish(key, event.data);
            }
          }

          eventBoundary = buffer.indexOf('\n\n');
        }
      }
    } catch (error) {
      if (error instanceof Error && error.name !== 'AbortError') {
        console.error('Error in SSE subscription:', error);
      }
    } finally {
      console.log(`SSE connection for ${key} has ended`);
    }
  })();

  const idleCheckInterval = setInterval(() => {
    if (activeSubscriptions.has(key)) {
      const { lastActivity, controller } = activeSubscriptions.get(key)!;
      const inactiveTime = Date.now() - lastActivity;

      if (inactiveTime > 5 * 60 * 1000) {
        controller.abort();
        activeSubscriptions.delete(key);
        clearInterval(idleCheckInterval);
        window.removeEventListener('beforeunload', unloadListener);
      }
    } else {
      clearInterval(idleCheckInterval);
    }
  }, 60000);

  return {
    cancel: () => {
      PubSub.publish(`cleanup-${key}`, null);
      clearInterval(idleCheckInterval);
    },
  };
}
