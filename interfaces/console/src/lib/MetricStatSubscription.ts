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
  siteId?: string;
  type: Stats_Type;
}

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

export default async function MetricStatSubscription({
  url,
  key,
  from,
  type,
  userId,
  siteId,
  nodeId,
  orgName,
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

  const res = await fetch(
    `${url}/graphql?query=subscription+MetricStatSub%28%24data%3ASubMetricsStatInput%21%29%7BgetMetricStatSub%28data%3A%24data%29%7Bmsg+nodeId+success+type+value%7D%7D&variables=%7B%22data%22%3A%7B%22nodeId%22%3A%22${nodeId}%22%2C%22orgName%22%3A%22${orgName}%22%2C%22type%22%3A%22${type}%22%2C%22userId%22%3A%22${userId}%22%2C%22siteId%22%3A%22${siteId}%22%2C%22from%22%3A${from}%7D%7D&operationName=MetricStatSub&extensions=%7B%7D`,
    { ...requestOptions, signal },
  ).catch((error) => {
    if (error.name === 'AbortError') {
      console.log('Fetch aborted');
    } else {
      console.error('Fetch error:', error);
    }
  });

  if (!res || !res.ok) {
    console.error('Network response was not ok');
    return;
  }

  const reader = res?.body?.getReader();
  const decoder = new TextDecoder('utf-8');
  let buffer = '';

  while (true) {
    const { value, done } = (await reader?.read()) || {};

    if (done) {
      console.log('Stream complete');
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
}
