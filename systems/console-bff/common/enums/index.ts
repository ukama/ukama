/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import "reflect-metadata";
import { registerEnumType } from "type-graphql";

export enum API_METHOD_TYPE {
  GET = "get",
  POST = "post",
  PUT = "put",
  DELETE = "delete",
  PATCH = "patch",
}
export enum NODE_STATE {
  Unknown = "Unknown",
  Configured = "Configured",
  Operational = "Operational",
  Faulty = "Faulty",
}
registerEnumType(NODE_STATE, {
  name: "NodeStateEnum",
  description: "Node state enums",
});

export enum NODE_CONNECTIVITY {
  Unknown = "Unkown",
  Offline = "Offline",
  Online = "Online",
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
  unknown = "unknown",
  test = "test",
  operator_data = "operator_data",
  ukama_data = "ukama_data",
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
  HOME = "HOME",
  NODE_HEALTH = "NODE_HEALTH",
  SUBSCRIBERS = "SUBSCRIBERS",
  DATA_USAGE = "DATA_USAGE",
  NETWORK_CELLULAR = "NETWORK_CELLULAR",
  NETWORK_BACKHAUL = "NETWORK_BACKHAUL",
  RESOURCES = "RESOURCES",
  RADIO = "RADIO",
  SOLAR = "SOLAR",
  BATTERY = "BATTERY",
  CONTROLLER = "CONTROLLER",
  MAIN_BACKHAUL = "MAIN_BACKHAUL",
  SWITCH = "SWITCH",
  SITE = "SITE",
}
registerEnumType(GRAPHS_TYPE, {
  name: "GRAPHS_TYPE",
});

export enum STATS_TYPE {
  ALL_NODE = "ALL_NODE",
  HOME = "HOME",
  OVERVIEW = "OVERVIEW",
  DATA_USAGE = "DATA_USAGE",
  NETWORK = "NETWORK",
  RESOURCES = "RESOURCES",
  RADIO = "RADIO",
  SITE = "SITE",
  BATTERY = "BATTERY",
  MAIN_BACKHAUL = "MAIN_BACKHAUL",
}
registerEnumType(STATS_TYPE, {
  name: "STATS_TYPE",
});

export enum ROLE_TYPE {
  ROLE_INVALID = "ROLE_INVALID",
  ROLE_OWNER = "ROLE_OWNER",
  ROLE_ADMIN = "ROLE_ADMIN",
  ROLE_NETWORK_OWNER = "ROLE_NETWORK_OWNER",
  ROLE_VENDOR = "ROLE_VENDOR",
  ROLE_USER = "ROLE_USER",
}
registerEnumType(ROLE_TYPE, {
  name: "ROLE_TYPE",
});

export enum SIM_STATUS {
  ALL = "ALL",
  ASSIGNED = "ASSIGNED",
  UNASSIGNED = "UNASSIGNED",
}
registerEnumType(SIM_STATUS, {
  name: "SIM_STATUS",
});

export enum NOTIFICATION_TYPE {
  TYPE_INVALID = "TYPE_INVALID",
  TYPE_INFO = "TYPE_INFO",
  TYPE_WARNING = "TYPE_WARNING",
  TYPE_ERROR = "TYPE_ERROR",
  TYPE_CRITICAL = "TYPE_CRITICAL",
  TYPE_ACTIONABLE_INFO = "TYPE_ACTIONABLE_INFO",
  TYPE_ACTIONABLE_WARNING = "TYPE_ACTIONABLE_WARNING",
  TYPE_ACTIONABLE_ERROR = "TYPE_ACTIONABLE_ERROR",
  TYPE_ACTIONABLE_CRITICAL = "TYPE_ACTIONABLE_CRITICAL",
}
registerEnumType(NOTIFICATION_TYPE, {
  name: "NOTIFICATION_TYPE",
});

export enum NOTIFICATION_SCOPE {
  SCOPE_INVALID = "SCOPE_INVALID",
  SCOPE_OWNER = "SCOPE_OWNER",
  SCOPE_ORG = "SCOPE_ORG",
  SCOPE_NETWORKS = "SCOPE_NETWORKS",
  SCOPE_NETWORK = "SCOPE_NETWORK",
  SCOPE_SITES = "SCOPE_SITES",
  SCOPE_SITE = "SCOPE_SITE",
  SCOPE_SUBSCRIBERS = "SCOPE_SUBSCRIBERS",
  SCOPE_SUBSCRIBER = "SCOPE_SUBSCRIBER",
  SCOPE_USERS = "SCOPE_USERS",
  SCOPE_USER = "SCOPE_USER",
  SCOPE_NODE = "SCOPE_NODE",
}
registerEnumType(NOTIFICATION_SCOPE, {
  name: "NOTIFICATION_SCOPE",
});

export enum INVITATION_STATUS {
  INVITE_PENDING = "INVITE_PENDING",
  INVITE_ACCEPTED = "INVITE_ACCEPTED",
  INVITE_DECLINED = "INVITE_DECLINED",
}
registerEnumType(INVITATION_STATUS, {
  name: "INVITATION_STATUS",
});
export enum PAYMENT_ITEM_TYPE {
  UNKNOWN = "unknown",
  PACKAGE = "package",
  INVOICE = "invoice",
}
registerEnumType(PAYMENT_ITEM_TYPE, {
  name: "PAYMENT_ITEM_TYPE",
});

export enum COMPONENT_TYPE {
  all = "all",
  access = "access",
  backhaul = "backhaul",
  power = "power",
  switch = "switch",
  spectrum = "spectrum",
}
registerEnumType(COMPONENT_TYPE, {
  name: "COMPONENT_TYPE",
});

export const NotificationScopeEnumValue = (e: number) => {
  switch (e) {
    case 0:
      return NOTIFICATION_SCOPE.SCOPE_INVALID;
    case 1:
      return NOTIFICATION_SCOPE.SCOPE_OWNER;
    case 2:
      return NOTIFICATION_SCOPE.SCOPE_ORG;
    case 3:
      return NOTIFICATION_SCOPE.SCOPE_NETWORKS;
    case 4:
      return NOTIFICATION_SCOPE.SCOPE_NETWORK;
    case 5:
      return NOTIFICATION_SCOPE.SCOPE_SITES;
    case 6:
      return NOTIFICATION_SCOPE.SCOPE_SITE;
    case 7:
      return NOTIFICATION_SCOPE.SCOPE_SUBSCRIBERS;
    case 8:
      return NOTIFICATION_SCOPE.SCOPE_SUBSCRIBER;
    case 9:
      return NOTIFICATION_SCOPE.SCOPE_USERS;
    case 10:
      return NOTIFICATION_SCOPE.SCOPE_USER;
    case 11:
      return NOTIFICATION_SCOPE.SCOPE_NODE;
    default:
      return NOTIFICATION_SCOPE.SCOPE_INVALID;
  }
};

export const NotificationTypeEnumValue = (e: number) => {
  switch (e) {
    case 0:
      return NOTIFICATION_TYPE.TYPE_INVALID;
    case 1:
      return NOTIFICATION_TYPE.TYPE_INFO;
    case 2:
      return NOTIFICATION_TYPE.TYPE_WARNING;
    case 3:
      return NOTIFICATION_TYPE.TYPE_ERROR;
    case 4:
      return NOTIFICATION_TYPE.TYPE_CRITICAL;
    case 5:
      return NOTIFICATION_TYPE.TYPE_ACTIONABLE_INFO;
    case 6:
      return NOTIFICATION_TYPE.TYPE_ACTIONABLE_WARNING;
    case 7:
      return NOTIFICATION_TYPE.TYPE_ACTIONABLE_ERROR;
    case 8:
      return NOTIFICATION_TYPE.TYPE_ACTIONABLE_CRITICAL;
    default:
      return NOTIFICATION_TYPE.TYPE_INVALID;
  }
};
