/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Stats_Type } from '@/client/graphql/generated/subscriptions';
import PubSub from 'pubsub-js';

interface IMetricStatSubscription {
  url: string;
  key: string;
  from: number;
  nodeId?: string;
  userId: string;
  orgName: string;
  type: Stats_Type;
  networkId?: string;
}

function parseEvent(eventStr: any) {
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

  return event;
}

export default async function MetricStatSubscription({
  url,
  key,
  from,
  type,
  userId,
  orgName,
  nodeId = undefined,
  networkId = undefined,
}: IMetricStatSubscription) {
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

  let fullUrl = '';
  const baseParams: {
    nodeId?: string;
    orgName: string;
    type: Stats_Type;
    userId: string;
    from: number;
    siteId?: string;
    networkId?: string;
  } = {
    orgName,
    type,
    userId,
    from,
  };

  if (nodeId) {
    baseParams.nodeId = nodeId;
  }
  if (networkId) {
    baseParams.networkId = networkId;
  }

  const queryParams = new URLSearchParams({
    query: 'subscription MetricStatSub($data:SubMetricsStatInput!){getMetricStatSub(data:$data){msg nodeId success type value networkId packageId dataPlanId}}',
    operationName: 'MetricStatSub',
    extensions: '{}',
  });

  const variables = { data: baseParams };

  queryParams.append('variables', JSON.stringify(variables));

  fullUrl = `${url}/graphql?${queryParams.toString()}`;

  const res = await fetch(fullUrl, { ...requestOptions, signal }).catch(
    (error) => {
      if (error.name === 'AbortError') {
        console.log('Fetch aborted');
      } else {
        console.error('Fetch error:', error);
      }
    },
  );

  if (!res || !res.ok) {
    console.error('Network response was not ok');
    return controller;
  }

  const reader = res?.body?.getReader();
  const decoder = new TextDecoder('utf-8');
  let buffer = '';

  const processStream = async () => {
    try {
      while (true) {
        const { value, done } = (await reader?.read()) || {};

        if (done) {
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
      }
    } catch (error: any) {
      if (error.name === 'AbortError') {
      } else {
        console.error('Stream error:', error);
      }
    } finally {
      try {
        if (reader && !reader.closed) {
          await reader.cancel();
        }
      } catch (error: any) {
        if (error.name !== 'AbortError') {
          console.error('Error canceling reader:', error);
        }
      }
    }
  };

  processStream().catch((error) => {
    console.error('Error in processStream:', error);
  });

  return controller;
}
