/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { kvsLocalStorage } from "@kvs/node-localstorage";

import { STORAGE_KEY } from "../configs";

const storeInStorage = async (key: string, value: any) => {
  const storageKeyValue = await kvsLocalStorage({
    name: STORAGE_KEY,
    version: 1,
  });
  await storageKeyValue.set(key, value);
};

const retriveFromStorage = async (key: string): Promise<any> => {
  const storageKeyValue = await kvsLocalStorage({
    name: STORAGE_KEY,
    version: 1,
  });
  return await storageKeyValue.get(key);
};

const removeKeyFromStorage = async (key: string) => {
  const storageKeyValue = await kvsLocalStorage({
    name: STORAGE_KEY,
    version: 1,
  });
  await storageKeyValue.delete(key);
};

export { removeKeyFromStorage, retriveFromStorage, storeInStorage };
