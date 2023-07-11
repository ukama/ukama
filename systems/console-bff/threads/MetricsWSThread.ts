const { kvsLocalStorage } = require("@kvs/node-localstorage");
const { workerData, isMainThread, parentPort } = require("worker_threads");

const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

const storeInStorage = async (key: string, value: any, storageKey: string) => {
  const storageKeyValue = await kvsLocalStorage({
    name: storageKey,
    version: 1,
  });
  await storageKeyValue.set(key, value);
};

const retriveFromStorage = async (
  key: string,
  storageKey: string
): Promise<any> => {
  const storageKeyValue = await kvsLocalStorage({
    name: storageKey,
    version: 1,
  });
  return await storageKeyValue.get(key);
};

const removeKeyFromStorage = async (key: string, storageKey: string) => {
  const storageKeyValue = await kvsLocalStorage({
    name: storageKey,
    version: 1,
  });
  await storageKeyValue.delete(key);
};

const runWorker = async () => {
  if (!isMainThread) {
    const WebSocket = require("ws");
    const { url, orgId, userId, type, key: storageKey } = workerData;
    const ws = new WebSocket(url);

    ws.on("error", e =>
      parentPort.postMessage({ isError: true, message: e.message, data: null })
    );

    ws.on("open", async function open() {
      await storeInStorage(
        `${orgId}/${userId}/${type}`,
        getTimestampCount("0"),
        storageKey
      );
    });

    ws.on("message", async function message(data) {
      const value = await retriveFromStorage(
        `${orgId}/${userId}/${type}`,
        storageKey
      );
      let occurance = value ? parseInt(value.split("-")[1]) : 0;
      occurance += 1;
      await storeInStorage(
        `${orgId}/${userId}/${type}`,
        getTimestampCount(`${occurance}`),
        storageKey
      );
      if (occurance === 9) {
        ws.close();
        ws.terminate();
        removeKeyFromStorage(`${orgId}/${userId}/${type}`, storageKey);
        process.exit(0);
      }
      parentPort.postMessage({
        isError: false,
        message: "",
        data: data.toString(),
      });
    });
  }
};

runWorker();
