import { registerEnumType } from "type-graphql";

export const PORT = process.env.PORT;
export const DEV_URL = process.env.API_URL;
export const BASE_URL = `http://localhost:${PORT}`;
export const HEADER = {
    headers: {
        cookie: "ukama_session=test",
    },
};

export enum CONNECTED_USER_TYPE {
    RESIDENTS = "RESIDENTS",
    GUESTS = "GUESTS",
}
registerEnumType(CONNECTED_USER_TYPE, {
    name: "CONNECTED_USER_TYPE",
});

export enum TIME_FILTER {
    TODAY = "TODAY",
    WEEK = "WEEK",
    MONTH = "MONTH",
    TOTAL = "TOTAL",
}
registerEnumType(TIME_FILTER, {
    name: "TIME_FILTER",
});

export enum DATA_BILL_FILTER {
    CURRENT = "CURRENT",
    JANUARY = "JANUARY",
    FEBRURAY = "FEBRURAY",
    MARCH = "MARCH",
    APRIL = "APRIL",
    MAY = "MAY",
    JUNE = "JUNE",
    JULY = "JULY",
    AUGUST = "AUGUST",
    SEPTEMBER = "SEPTEMBER",
    OCTOBER = "OCTOBER",
    NOVERMBER = "NOVERMBER",
    DECEMBER = "DECEMBER",
}
registerEnumType(DATA_BILL_FILTER, {
    name: "DATA_BILL_FILTER",
});

export enum ALERT_TYPE {
    INFO = "INFO",
    WARNING = "WARNING",
    ERROR = "ERROR",
}
registerEnumType(ALERT_TYPE, {
    name: "ALERT_TYPE",
});

export enum API_METHOD_TYPE {
    GET = "get",
    POST = "post",
    PUT = "put",
    DELETE = "delete",
    PATCH = "patch",
}
registerEnumType(API_METHOD_TYPE, {
    name: "API_METHOD_TYPE",
});

export enum GET_STATUS_TYPE {
    UNKNOWN = "UNKNOWN",
    ACTIVE = "ACTIVE",
    INACTIVE = "INACTIVE",
}
registerEnumType(GET_STATUS_TYPE, {
    name: "GET_USER_STATUS_TYPE",
});

export enum GET_USER_TYPE {
    ALL = "ALL",
    RESIDENT = "RESIDENT",
    VISITOR = "VISITOR",
    HOME = "HOME",
    GUEST = "GUEST",
}
registerEnumType(GET_USER_TYPE, {
    name: "GET_USER_TYPE",
});
export enum NETWORK_STATUS {
    ONLINE = "ONLINE",
    BEING_CONFIGURED = "BEING_CONFIGURED",
}
registerEnumType(NETWORK_STATUS, {
    name: "NETWORK_STATUS",
});

export enum NETWORK_TYPE {
    PUBLIC = "PUBLIC",
    PRIVATE = "PRIVATE",
}
registerEnumType(NETWORK_TYPE, {
    name: "NETWORK_TYPE",
});

export enum ORG_NODE_STATE {
    ONBOARDED = "ONBOARDED",
    PENDING = "PENDING",
    UNDEFINED = "UNDEFINED",
    ERROR = "UNDEFINED",
}
registerEnumType(ORG_NODE_STATE, {
    name: "ORG_NODE_STATE",
});
export enum GRAPH_FILTER {
    DAY = "DAY",
    WEEK = "WEEK",
    MONTH = "MONTH",
}
registerEnumType(GRAPH_FILTER, {
    name: "GRAPH_FILTER",
});
export enum NODE_TYPE {
    TOWER = "TOWER",
    AMPLIFIER = "AMPLIFIER",
    HOME = "HOME",
}
registerEnumType(NODE_TYPE, {
    name: "NODE_TYPE",
});
export enum GRAPHS_TAB {
    HOME = "HOME",
    RADIO = "RADIO",
    NETWORK = "NETWORK",
    OVERVIEW = "OVERVIEW",
    RESOURCES = "RESOURCES",
    NODE_STATUS = "NODE_STATUS",
}
registerEnumType(GRAPHS_TAB, {
    name: "GRAPHS_TAB",
});
