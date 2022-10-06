// eslint-disable-next-line @typescript-eslint/no-var-requires
const { workerData, isMainThread, parentPort } = require("worker_threads");

const oneSecSleep = (t = 1000) => new Promise(res => setTimeout(res, t));

const runWorker = async () => {
    if (!isMainThread) {
        for (let i = 0; i < workerData.length; i++) {
            const metric = [];
            for (const element of workerData.response) {
                metric.push({
                    next:
                        element.data[i] && i === workerData.length - 1
                            ? true
                            : false,
                    type: element.type,
                    name: element.name,
                    data: element.data[i] ? [element.data[i]] : [],
                });
            }
            await oneSecSleep();
            parentPort.postMessage({ metric });
        }
    }
};

runWorker();
