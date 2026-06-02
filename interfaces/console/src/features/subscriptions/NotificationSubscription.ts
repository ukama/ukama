/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import PubSub from 'pubsub-js';

interface SSEEvent {
  data?: string;
  id?: string;
  event?: string;
}

function parseEvent(eventStr: string): SSEEvent {
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

  return event;
}

export default async function ServerNotificationSubscription(
  url: string,
  key: string,
  role: string,
  orgId: string,
  userId: string,
  orgName: string,
  networkId: string,
  startTimestamp: string,
  signal?: AbortSignal,
) {
  const myHeaders = new Headers();
  myHeaders.append('Cache-Control', 'no-cache');
  myHeaders.append('Connection', 'keep-alive');
  myHeaders.append('Pragma', 'no-cache');
  myHeaders.append('Sec-Fetch-Dest', 'empty');
  myHeaders.append('Sec-Fetch-Mode', 'cors');
  myHeaders.append('Sec-Fetch-Site', 'same-origin');
  myHeaders.append('accept', 'text/event-stream');

  const query =
    `subscription NotificationSubscription{notificationSubscription(` +
    `networkId:"${encodeURIComponent(networkId)}" ` +
    `orgId:"${encodeURIComponent(orgId)}" ` +
    `orgName:"${encodeURIComponent(orgName)}" ` +
    `role:"${encodeURIComponent(role)}" ` +
    `startTimestamp:"${encodeURIComponent(startTimestamp)}" ` +
    `subscriberId:"" ` +
    `userId:"${encodeURIComponent(userId)}"` +
    `){id title description createdAt isRead scope type eventKey resourceId redirect{action title}}}`;

  const sseUrl = new URL(`${url}/graphql`);
  sseUrl.searchParams.set('query', query);
  sseUrl.searchParams.set('operationName', 'NotificationSubscription');
  sseUrl.searchParams.set('extensions', '{}');

  let res: Response;
  try {
    res = await fetch(sseUrl.toString(), { method: 'GET', headers: myHeaders, signal });
  } catch (err) {
    if ((err as Error).name !== 'AbortError') {
      console.error('[NotificationSubscription] fetch error:', err);
    }
    return;
  }

  const reader = res?.body?.getReader();
  if (!reader) return;

  const decoder = new TextDecoder('utf-8');
  let buffer = '';

  try {
    while (true) {
      if (signal?.aborted) break;

      const { value, done } = await reader.read();

      if (done) break;

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
  } catch (err) {
    if ((err as Error).name !== 'AbortError') {
      console.error('[NotificationSubscription] stream error:', err);
    }
  } finally {
    reader.cancel();
  }
}
