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

function parseEvent(eventStr: string) {
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

export default async function MetricStatBySiteSubscription({
  url,
  key,
  from,
  type,
  userId,
  siteIds,
  nodeIds,
  orgName,
}: IMetricStatBySiteSubscription) {
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

  const siteIdsParam = encodeURIComponent(JSON.stringify(siteIds));

  let fullUrl = `${url}/graphql?query=subscription+SiteMetricStatSub%28%24data%3ASubSiteMetricsStatInput%21%29%7BgetSiteMetricStatSub%28data%3A%24data%29%7Bmsg+siteId+nodeId+success+type+value%7D%7D&variables=%7B%22data%22%3A%7B%22siteIds%22%3A${siteIdsParam}%2C%22orgName%22%3A%22${orgName}%22%2C%22type%22%3A%22${type}%22%2C%22userId%22%3A%22${userId}%22%2C%22from%22%3A${from}`;

  if (nodeIds && nodeIds.length > 0) {
    const nodeIdsParam = encodeURIComponent(JSON.stringify(nodeIds));
    fullUrl += `%2C%22nodeIds%22%3A${nodeIdsParam}`;
  }

  fullUrl += `%7D%7D&operationName=SiteMetricStatSub&extensions=%7B%7D`;

  try {
    const res = await fetch(fullUrl, { ...requestOptions, signal }).catch(
      (error) => {
        if (error.name === 'AbortError') {
          console.log('Fetch aborted');
        } else {
          console.error('Fetch error:', error);
        }
        return null;
      },
    );

    if (!res || !res.ok) {
      console.error('Network response was not ok');
      return;
    }

    const reader = res.body?.getReader();
    if (!reader) {
      console.error('Could not get reader from response body');
      return;
    }

    const decoder = new TextDecoder('utf-8');
    let buffer = '';

    while (true) {
      const { value, done } = await reader.read();

      if (done) {
        console.log('Stream complete');
        break;
      }

      buffer += decoder.decode(value, { stream: true });

      let eventBoundary = buffer.indexOf('\n\n');
      while (eventBoundary !== -1) {
        const eventStr = buffer.slice(0, eventBoundary).trim();
        buffer = buffer.slice(eventBoundary + 2);

        if (eventStr !== ':') {
          const pevent = parseEvent(eventStr);

          if (pevent.data) {
            PubSub.publish(key, pevent.data);
          }
        }

        eventBoundary = buffer.indexOf('\n\n');
      }
    }
  } catch (error) {
    console.error('Error in subscription:', error);
  } finally {
    // Clean up
    controller.abort();
  }
}
