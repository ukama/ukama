import { registerEnumType } from "type-graphql";

export const NODE_ENV = "development";
export const PORT = "8080";
export const BASE_URL = `http://localhost:${PORT}`;

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
