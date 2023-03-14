import { Meta } from "../common/types";
import { GRAPHS_TAB, GRAPH_FILTER, NODE_TYPE } from "../constants";
import { AddNodeDto, NodeObj, LinkNodes } from "../modules/node/types";

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
                    "uptimetrx",
                    "temperaturetrx",
                    "temperaturerfe",
                    "subscribersactive",
                    "subscribersattached",
                ];
            else if (nodeType === NODE_TYPE.AMPLIFIER)
                return ["temperaturectl", "temperaturerfe"];
            else
                return [
                    "uptimetrx",
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
            return ["uptimetrx"];

        case GRAPHS_TAB.NODE_STATUS:
            if (nodeType === NODE_TYPE.HOME) return ["uptimetrx"];
            else if (nodeType === NODE_TYPE.AMPLIFIER) return ["uptimectl"];
            else return ["uptimetrx"];
    }
};

export const getMetricTitleByType = (type: string): string => {
    switch (type) {
        case "uptimetrx":
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

export const converCookieToObj = (cookie: string) => {
    if (cookie) {
        return cookie.split(";").reduce((res, c) => {
            const [key] = c.trim().split("=").map(decodeURIComponent);
            const val = c.slice(c.indexOf("=") + 1);
            try {
                return Object.assign(res, { [key]: JSON.parse(val) });
            } catch (e) {
                return Object.assign(res, { [key]: val });
            }
        }, {});
    }
    return null;
};

export const isTowerNode = (nodeId: string): boolean =>
    nodeId.includes("tnode");

export const getTowerNode = (payload: AddNodeDto): NodeObj => {
    if (isTowerNode(payload.nodeId))
        return {
            name: payload.name,
            state: payload.state,
            nodeId: payload.nodeId,
        };

    if (payload.attached)
        for (const node of payload.attached) {
            if (isTowerNode(node.nodeId)) return node;
        }
    return { name: "", nodeId: "", state: "" };
};

export const getNodes = (payload: AddNodeDto): NodeObj[] => {
    const nodes: NodeObj[] = [];
    if (!isTowerNode(payload.nodeId)) {
        nodes.push({
            name: payload.name,
            state: payload.state,
            nodeId: payload.nodeId,
        });
    }
    if (payload.attached)
        for (const node of payload.attached) {
            if (!isTowerNode(node.nodeId))
                nodes.push({
                    name: node.name,
                    state: payload.state,
                    nodeId: node.nodeId,
                });
        }
    return nodes;
};
export const linkNodes = (nodes: NodeObj[], rootNodeId: string): LinkNodes => {
    const nodesLinkingObj: LinkNodes = {
        nodeId: rootNodeId,
        attachedNodeIds: [],
    };
    for (const node of nodes) {
        nodesLinkingObj.attachedNodeIds?.push(node.nodeId);
    }
    return nodesLinkingObj;
};

export const getStripeIdByUserId = (uid: string): string => {
    switch (uid) {
        case "d0a36c51-6a66-4187-b786-72a9e09bf7a4":
            return "cus_MFTZKUVOGtI2fU";
        default:
            return "cus_MFTZKUVOGtI2fU";
    }
};
