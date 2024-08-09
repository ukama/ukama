/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { open } from "lmdb";
import { isMainThread, parentPort, workerData } from "worker_threads";
import WebSocket from "ws";

const TIMESTAMP_DIVISOR = 1000;
const MAX_OCCURRENCE = 5;

const getTimestampCount = count =>
  `${Math.floor(Date.now() / TIMESTAMP_DIVISOR)}-${count}`;

const openStore = () => {
  const path = process.env.STORAGE_KEY;
  if (!path) {
    throw new Error("STORAGE_KEY environment variable is not set");
  }
  return open({
    path,
    compression: true,
  });
};

const addInStore = async (store, key, value) => {
  try {
    await store.put(key, value);
  } catch (error) {
    console.error(`Failed to add in store: ${error.message}`);
  }
};

const getFromStore = async (store, key) => {
  try {
    return await store.get(key);
  } catch (error) {
    console.error(`Failed to get from store: ${error.message}`);
    return null;
  }
};

const constructParams = ({ orgId, userId, networkId, subscriberId }) => {
  const params = new URLSearchParams();
  if (orgId) params.append("org_id", orgId);
  if (userId) params.append("user_id", userId);
  if (networkId) params.append("network_id", networkId);
  if (subscriberId) params.append("subscriber_id", subscriberId);
  return params.toString();
};

const runWorker = async () => {
  if (!isMainThread) {
    const store = openStore();
    const { url, key, scopes } = workerData;
    const params = constructParams(workerData);
    const reqUrl = `${url}?${params}&${scopes}`;
    console.log("Opeing WebSocket connection to: ", reqUrl);
    const ws = new WebSocket(reqUrl);

    ws.on("error", e => {
      console.error("WebSocket error: ", e.message);
      parentPort.postMessage({ isError: true, message: e.message, data: null });
    });

    ws.on("open", async () => {
      console.error("WebSocket opened");
      await addInStore(store, key, getTimestampCount("0"));
    });

    ws.on("message", async data => {
      console.error("WebSocket message", data.toString());
      const value = await getFromStore(store, key);
      let occurrence = value ? parseInt(value.split("-")[1]) : 0;
      occurrence += 1;
      await addInStore(store, key, getTimestampCount(`${occurrence}`));
      if (occurrence === MAX_OCCURRENCE) {
        ws.terminate();
        ws.close();
      }

      parentPort.postMessage({
        isError: false,
        message: "success",
        data: data.toString(),
      });
    });
  }
};

runWorker();
