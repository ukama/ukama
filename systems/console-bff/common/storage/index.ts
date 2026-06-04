/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RootDatabase, open } from "lmdb";

import { STORAGE_KEY } from "../configs";

/**
 * Stored values are wrapped so a per-key expiry can travel with them. Entries
 * written before this wrapper existed are read back as-is (legacy fallback).
 */
interface TTLEntry {
  __ttl: true;
  /** Absolute expiry in epoch ms, or null for no expiry. */
  exp: number | null;
  value: unknown;
}

const isTTLEntry = (v: unknown): v is TTLEntry =>
  typeof v === "object" && v !== null && (v as TTLEntry).__ttl === true;

const openStore = (): RootDatabase => {
  return open({
    path: STORAGE_KEY,
    compression: true,
    maxReaders: 1024,
  });
};

/**
 * Stores a value, optionally with a time-to-live. When ttlSeconds is omitted
 * the entry never expires (previous behaviour).
 */
const addInStore = async (
  store: RootDatabase,
  key: string,
  value: any,
  ttlSeconds?: number
): Promise<boolean> => {
  const entry: TTLEntry = {
    __ttl: true,
    exp: ttlSeconds ? Date.now() + ttlSeconds * 1000 : null,
    value,
  };
  return await store.put(key, entry);
};

/**
 * Reads a value, transparently unwrapping the TTL envelope. Expired entries
 * are removed lazily on access and reported as missing (undefined).
 */
const getFromStore = async (store: RootDatabase, key: string) => {
  const raw = await store.get(key);
  if (raw === undefined) return undefined;
  if (!isTTLEntry(raw)) return raw; // legacy un-wrapped value
  if (raw.exp !== null && Date.now() >= raw.exp) {
    await store.remove(key);
    return undefined;
  }
  return raw.value;
};

const removeFromStore = async (
  store: RootDatabase,
  key: string
): Promise<boolean> => {
  return await store.remove(key);
};

/** Flushes and closes the LMDB environment (used during graceful shutdown). */
const closeStore = async (store: RootDatabase): Promise<void> => {
  await store.close();
};

export { addInStore, closeStore, getFromStore, openStore, removeFromStore };
