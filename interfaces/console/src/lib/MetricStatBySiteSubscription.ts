/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Stats_Type } from '@/client/graphql/generated/subscriptions';
import PubSub from 'pubsub-js';

interface IMetricBySiteStatSubscription {
  url: string;
  key: string;
  from: number;
  siteId: string;
  userId: string;
  orgName: string;
  type: Stats_Type;
  nodeIds?: string[];
}

const activeSubscriptions = new Map();

function parseEvent(eventStr: any) {
  const event: any = {};
  const lines = eventStr.split('\n');

  for (let line of lines) {
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

export function cancelSubscription(key: string) {
  const controller = activeSubscriptions.get(key);
  if (controller) {
    controller.abort();
    activeSubscriptions.delete(key);
    console.log(`Cancelled subscription for key: ${key}`);
    return true;
  }
  return false;
}

export default async function MetricStatBySiteSubscription({
  url,
  key,
  from,
  type,
  userId,
  siteId,
  nodeIds,
  orgName,
}: IMetricBySiteStatSubscription) {
  cancelSubscription(key);

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
  };

  const controller = new AbortController();
  const signal = controller.signal;

  activeSubscriptions.set(key, controller);

  type SubSiteMetricsStatInput = {
    siteId: string;
    orgName: string;
    type: Stats_Type;
    userId: string;
    from: number;
    nodeIds?: string[];
  };

  const data: SubSiteMetricsStatInput = {
    siteId,
    orgName,
    type,
    userId,
    from,
  };

  if (nodeIds && nodeIds.length > 0) {
    data.nodeIds = nodeIds;
  }

  const query =
    'subscription+SiteMetricStatSub%28%24data%3ASubSiteMetricsStatInput%21%29%7BgetSiteMetricStatSub%28data%3A%24data%29%7Bmsg+siteId+nodeId+success+type+value%7D%7D';

  const variables = encodeURIComponent(JSON.stringify({ data }));

  const fullUrl = `${url}/graphql?query=${query}&variables=${variables}&operationName=SiteMetricStatSub&extensions=%7B%7D`;

  try {
    const res = await fetch(fullUrl, { ...requestOptions, signal });

    if (!res || !res.body) {
      console.error('Error: No response or response body');
      activeSubscriptions.delete(key);
      return;
    }

    const reader = res.body.getReader();
    const decoder = new TextDecoder('utf-8');
    let buffer = '';

    while (true) {
      try {
        const { value, done } = await reader.read();

        if (done) {
          console.log(`Stream complete for key: ${key}`);
          break;
        }

        buffer += decoder.decode(value, { stream: true });

        let eventBoundary = buffer.indexOf('\n\n');

        while (eventBoundary !== -1) {
          const eventStr = buffer.slice(0, eventBoundary).trim();
          buffer = buffer.slice(eventBoundary + 2);
          const pevent = parseEvent(eventStr);
          if (pevent.data) {
            PubSub.publish(key, pevent.data);
          }
          eventBoundary = buffer.indexOf('\n\n');
        }
      } catch (error) {
        if (error instanceof Error && error.name === 'AbortError') {
          console.log(`Subscription for ${key} was aborted`);
        } else {
          console.error(`Error in subscription for ${key}:`, error);
        }
        break;
      }
    }
  } catch (error) {
    if (
      error &&
      typeof error === 'object' &&
      'name' in error &&
      error.name === 'AbortError'
    ) {
      console.log(`Fetch for ${key} was aborted`);
    } else {
      console.error(`Fetch error for ${key}:`, error);
    }
  } finally {
    if (activeSubscriptions.get(key) === controller) {
      activeSubscriptions.delete(key);
    }
  }
}
