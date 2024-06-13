/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "reflect-metadata";
import { registerEnumType } from "type-graphql";

export enum COMPONENT_CATEGORY {
  ALL = 0,
  ACCESS = 1,
  BACKHAUL = 2,
  POWER = 3,
  SWITCH = 4,
}

registerEnumType(COMPONENT_CATEGORY, {
  name: "ComponentCategory",
  description: "Categories for components",
});

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

export enum GRAPHS_TYPE {
  NODE_HEALTH = "NODE_HEALTH",
  SUBSCRIBERS = "SUBSCRIBERS",
  NETWORK = "NETWORK",
  RESOURCES = "RESOURCES",
  RADIO = "RADIO",
}
registerEnumType(GRAPHS_TYPE, {
  name: "GRAPHS_TYPE",
});

export enum ROLE_TYPE {
  OWNER = "OWNER",
  ADMIN = "ADMIN",
  VENDOR = "VENDOR",
  USERS = "USERS",
}
registerEnumType(ROLE_TYPE, {
  name: "ROLE_TYPE",
});

export enum NOTIFICATION_TYPE {
  UNKNOWN = "UNKNOWN",
  INFO = "INFO",
  WARNING = "WARNING",
  ERROR = "ERROR",
}
registerEnumType(NOTIFICATION_TYPE, {
  name: "NOTIFICATION_TYPE",
});

export enum NOTIFICATION_SCOPE {
  ORG = "ORG",
  NETWORK = "NETWORK",
  SITE = "SITE",
  SUBSCRIBER = "SUBSCRIBER",
  USER = "USER",
  NODE = "NODE",
}
registerEnumType(NOTIFICATION_SCOPE, {
  name: "NOTIFICATION_SCOPE",
});
