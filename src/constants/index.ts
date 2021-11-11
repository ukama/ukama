import { registerEnumType } from "type-graphql";

export enum CONNECTED_USER_TYPE {
    RESIDENTS = "Residents",
    GUESTS = "Guests",
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
