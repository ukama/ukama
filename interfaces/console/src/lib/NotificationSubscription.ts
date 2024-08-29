/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import PubSub from 'pubsub-js';
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

export default async function ServerNotificationSubscription(
  key: string,
  role: string,
  orgId: string,
  userId: string,
  orgName: string,
  networkId: string,
  startTimestamp: string,
) {
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
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_METRIC_URL}/graphql?query=subscription+NotificationSubscription%28%24data%3AGetNotificationsInput%21%29%7BnotificationSubscription%28data%3A%24data%29%7BcreatedAt+description+id+isRead+scope+title+type%7D%7D&variables=%7B%22data%22%3A%7B%22networkId%22%3A%22${networkId}%22%2C%22orgId%22%3A%22${orgId}%22%2C%22orgName%22%3A%22${orgName}%22%2C%22role%22%3A%22${role}%22%2C%22subscriberId%22%3A%22${''}%22%2C%22userId%22%3A%22${userId}%22%2C%22startTimestamp%22%3A%22${startTimestamp}%22%7D%7D&operationName=NotificationSubscription&extensions=%7B%7D`,
    requestOptions,
  );

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
