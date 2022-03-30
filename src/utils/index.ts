import { Meta } from "../common/types";
import { GRAPHS_TAB, GRAPH_FILTER, NODE_TYPE } from "../constants";

export const getPaginatedOutput = (
    page: number,
    pageSize: number,
    count: number
): Meta => {
    return {
        count,
        page: page ? page : 1,
        size: pageSize ? pageSize : count,
        pages: pageSize ? Math.ceil(count / pageSize) : 1,
    };
};

export const getUniqueTimeStamp = (index?: number, length?: number): number =>
    new Date().valueOf() - (length ? length - 1000 * (index || 1) : 0);

export const getRecordsLengthByFilter = (
    filter: string | undefined
): number => {
    switch (filter) {
        case GRAPH_FILTER.WEEK:
            return 50;
        case GRAPH_FILTER.MONTH:
            return 100;
        default:
            return 10;
    }
};

export const oneSecSleep = (t = 1000): any =>
    new Promise(res => setTimeout(res, t));

export const getMetricsByTab = (
    nodeType: NODE_TYPE,
    tabType: GRAPHS_TAB
): string[] => {
    switch (tabType) {
        case GRAPHS_TAB.OVERVIEW:
            if (nodeType === NODE_TYPE.HOME)
                return [
                    "temperaturetrx",
                    "temperaturerfe",
                    "subscribersactive",
                    "subscribersattached",
                ];
            else if (nodeType === NODE_TYPE.AMPLIFIER)
                return ["temperaturectl", "temperaturerfe"];
            else
                return [
                    "uptime",
                    "temperaturetrx",
                    "temperaturecom",
                    "subscribersactive",
                    "subscribersattached",
                ];

        case GRAPHS_TAB.NETWORK:
            if (nodeType !== NODE_TYPE.AMPLIFIER)
                return [
                    "rrc",
                    "rlc",
                    "erab",
                    "throughputuplink",
                    "throughputdownlink",
                ];
            else return [];

        case GRAPHS_TAB.RESOURCES:
            if (nodeType === NODE_TYPE.HOME)
                return ["cputrxusage", "memorytrxused", "disktrxused"];
            else if (nodeType === NODE_TYPE.AMPLIFIER)
                return ["cpuctlused", "diskctlused", "memoryctlused"];
            else
                return [
                    "powerlevel",
                    "cputrxusage",
                    "cpucomusage",
                    "disktrxused",
                    "diskcomused",
                    "memorytrxused",
                    "memorycomused",
                ];

        case GRAPHS_TAB.RADIO:
            if (nodeType !== NODE_TYPE.HOME)
                return ["txpower", "rxpower", "papower"];
            else return [];

        case GRAPHS_TAB.HOME:
            return ["uptime"];
    }
};

export const getMetricTitleByType = (type: string): string => {
    switch (type) {
        case "uptime":
            return "Uptime";
        case "temperaturetrx":
            return "Temp. (TRX)";
        case "temperaturerfe":
            return "Temp. (RFE)";
        case "subscribersactive":
            return "Active";
        case "subscribersattached":
            return "Attached";
        case "temperaturectl":
            return "Temp. (CTL)";
        case "temperaturecom":
            return "Temp. (COM)";
        case "rrc":
            return "RRC CNX success";
        case "rlc":
            return "RLS  drop rate";
        case "erab":
            return "ERAB drop rate";
        case "throughputuplink":
            return "Throughput (U/L)";
        case "throughputdownlink":
            return "Throughput (D/L)";
        case "cputrxusage":
            return "CPU-TRX";
        case "memorytrxused":
            return "Memory-TRX";
        case "disktrxused":
            return "DISK-TRX";
        case "cpuctlused":
            return "CPU-CTL";
        case "diskctlused":
            return "DISK-CTL";
        case "memoryctlused":
            return "Memory-CTL";
        case "powerlevel":
            return "Power";
        case "cpucomusage":
            return "CPU-COM";
        case "diskcomused":
            return "DISK-COM";
        case "memorycomused":
            return "Memory-COM";
        case "txpower":
            return "TX Power";
        case "rxpower":
            return "RX Power";
        case "papower":
            return "PA Power";
        default:
            return "";
    }
};
