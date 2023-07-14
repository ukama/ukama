/* eslint-disable @typescript-eslint/no-var-requires */
const { workerData, isMainThread, parentPort } = require("worker_threads");
const axios = require("axios");

const runWorker = async () => {
  if (!isMainThread) {
    console.log("workerData", workerData);
    const config = {
      method: "post",
      maxBodyLength: Infinity,
      url: workerData.url,
      headers: {
        "Content-Type": "application/json",
      },
      data: workerData.data,
    };
    const res = await axios.request(config);
    if (res.status === 200) {
      parentPort.postMessage({
        isSuccess: true,
        status: res.status,
        message: res.statusText,
        data: res.data,
      });
    } else {
      parentPort.postMessage({
        isSuccess: false,
        status: res.status,
        message: res.statusText,
        data: null,
      });
    }
  }
};

runWorker();