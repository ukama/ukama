import { registerEnumType } from "type-graphql";

export const NODE_ENV = "development";
export const PORT = "8080";
export const BASE_URL = `http://localhost:${PORT}`;
export const DEV_URL = `https://api.dev.ukama.com`;
export const HEADER = {
    headers: {
        cookie: "ukama_session=MTYzOTM4MTgwNHxEdi1CQkFFQ180SUFBUkFCRUFBQVJfLUNBQUVHYzNSeWFXNW5EQThBRFhObGMzTnBiMjVmZEc5clpXNEdjM1J5YVc1bkRDSUFJSFpJWkhOWWFIRTBiRE51ZUUxSlNIVmFVRkYxVjA1T2JYbFFXVFpVYkUxenzvvP08sJTXG0bYYBeH4V9BRXfsDC_NJz4pzUavjLh5Tw==; csrf_token_0a4b9640203f0baf1ae8f999c46e57938950e3e29d87300080eb0fd3b129b396=7HK+rqcMH1L+t62RQWZMlYB+EndnEgOCeZUG1R22WhM=",
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
}
registerEnumType(API_METHOD_TYPE, {
    name: "API_METHOD_TYPE",
});

export enum GET_STATUS_TYPE {
    ACTIVE = "ACTIVE",
    INACTIVE = "INACTIVE",
}
registerEnumType(GET_STATUS_TYPE, {
    name: "GET_USER_STATUS_TYPE",
});

export enum DATA_PLAN_TYPE {
    NA = "NA",
    PAID = "PAID",
    UNPAID = "UNPAID",
}
registerEnumType(DATA_PLAN_TYPE, {
    name: "DATA_PLAN_TYPE",
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

export enum NODE_TYPE {
    HOME = "Home",
    WORK = "Work",
}
registerEnumType(NODE_TYPE, {
    name: "NODE_TYPE",
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
}
registerEnumType(ORG_NODE_STATE, {
    name: "ORG_NODE_STATE",
});
