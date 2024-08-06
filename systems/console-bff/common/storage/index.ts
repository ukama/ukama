/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RootDatabase, open } from "lmdb";

const openStore = (): RootDatabase => {
  return open({
    path: "ukama-systems-bff",
    compression: true,
  });
};

const addInStore = async (
  store: RootDatabase,
  key: string,
  value: any
): Promise<boolean> => {
  return await store.put(key, value);
};

const getFromStore = async (store: RootDatabase, key: string) => {
  return await store.get(key);
};

const removeFromStore = async (
  store: RootDatabase,
  key: string
): Promise<boolean> => {
  return await store.remove(key);
};

export { addInStore, getFromStore, openStore, removeFromStore };
