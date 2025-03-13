/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/* eslint-disable @typescript-eslint/no-var-requires */
import { open } from "lmdb";
import { isMainThread, parentPort, workerData } from "worker_threads";
import WebSocket from "ws";

const MAX_OCCURRENCE = 90;

const openStore = () => {
  const path = process.env.STORAGE_KEY;
  if (!path) {
    throw new Error("STORAGE_KEY environment variable is not set");
  }
  return open({
    path,
    compression: true,
    maxReaders: 1024,
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

const runWorker = async () => {
  if (!isMainThread) {
    const store = openStore();
    const { url, topic } = workerData;
    const ws = new WebSocket(url);
    ws.on("error", e => {
      store.close();
      console.error("WebSocket error: ", e.message);
      parentPort.postMessage({ isError: true, message: e.message, data: null });
    });

    ws.close = () => {
      store.close();
      console.error("WebSocket closed");
      ws.terminate();
    };

    ws.on("open", async function open() {
      console.error("WebSocket opened");
      await addInStore(store, topic, 0);
    });

    ws.on("message", async function message(data) {
      const value = await getFromStore(store, topic);
      let occurrence = parseInt(value ?? "0") + 1;
      await addInStore(store, topic, occurrence);
      if (occurrence > MAX_OCCURRENCE) {
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
