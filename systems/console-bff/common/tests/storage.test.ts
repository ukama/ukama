/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { type RootDatabase, open } from "lmdb";
import { tmpdir } from "os";
import { join } from "path";

import { addInStore, getFromStore, removeFromStore } from "../storage";

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

describe("storage TTL", () => {
  let store: RootDatabase;

  beforeAll(() => {
    store = open({ path: join(tmpdir(), `bff-store-test-${Date.now()}`) });
  });

  afterAll(async () => {
    await store.close();
  });

  it("stores and retrieves a value", async () => {
    await addInStore(store, "k1", "v1");
    expect(await getFromStore(store, "k1")).toBe("v1");
  });

  it("keeps values with no TTL indefinitely", async () => {
    await addInStore(store, "k2", { a: 1 });
    await sleep(60);
    expect(await getFromStore(store, "k2")).toEqual({ a: 1 });
  });

  it("expires a value once its TTL elapses", async () => {
    await addInStore(store, "k3", "v3", 0.05); // 50ms
    expect(await getFromStore(store, "k3")).toBe("v3");
    await sleep(90);
    expect(await getFromStore(store, "k3")).toBeUndefined();
  });

  it("reads back legacy un-wrapped values", async () => {
    await store.put("legacy", "raw-value");
    expect(await getFromStore(store, "legacy")).toBe("raw-value");
  });

  it("removes a value", async () => {
    await addInStore(store, "k4", "v4");
    await removeFromStore(store, "k4");
    expect(await getFromStore(store, "k4")).toBeUndefined();
  });
});
