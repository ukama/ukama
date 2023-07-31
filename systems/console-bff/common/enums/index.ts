import { registerEnumType } from "type-graphql";

export enum API_METHOD_TYPE {
  GET = "get",
  POST = "post",
  PUT = "put",
  DELETE = "delete",
  PATCH = "patch",
}
export enum NODE_STATUS {
  ACTIVE = "active",
  MAINTENANCE = "maintenance",
  FAULTY = "faulty",
  ONBOARDED = "onboarded",
  CONFIGURED = "configured",
  UNDEFINED = "undefined",
}
registerEnumType(NODE_STATUS, {
  name: "NodeStatusEnum",
  description: "Node status enums",
});

export enum NODE_CONNECTIVITY {
  UNKNOWN = "unkown",
  OFFLINE = "offline",
  ONLINE = "online",
}
registerEnumType(NODE_CONNECTIVITY, {
  name: "NodeConnectivityEnum",
  description: "Node connectivity enums",
});

export enum NODE_TYPE {
  tnode = "tnode",
  anode = "anode",
  hnode = "hnode",
}
registerEnumType(NODE_TYPE, {
  name: "NodeTypeEnum",
  description: "Node type enums",
});

export enum ALERT_TYPE {
  INFO = "INFO",
  WARNING = "WARNING",
  ERROR = "ERROR",
}
registerEnumType(ALERT_TYPE, {
  name: "AlertTypeEnum",
  description: "Alert type enums",
});

export enum NETWORK_STATUS {
  ONLINE = "ONLINE",
  DOWN = "DOWN",
  UNDEFINED = "UNDEFINED",
}
registerEnumType(NETWORK_STATUS, {
  name: "NETWORK_STATUS",
});

export enum SIM_TYPES {
  UNKNOWN = "unknown",
  TEST = "test",
  OPERATOR_DATA = "operator_data",
  UKAMA_DATA = "ukama_data",
}
registerEnumType(SIM_TYPES, {
  name: "SIM_TYPES",
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
