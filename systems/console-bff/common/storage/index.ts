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
