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

export enum NOTIFICATION_TYPE {
  NOTIF_INVALID = "NOTIF_INVALID",
  NOTIF_INFO = "NOTIF_INFO",
  NOTIF_WARNING = "NOTIF_WARNING",
  NOTIF_ERROR = "NOTIF_ERROR",
  NOTIF_CRITICAL = "NOTIF_CRITICAL",
  NOTIF_ACTIONABLE_INFO = "NOTIF_ACTIONABLE_INFO",
  NOTIF_ACTIONABLE_WARNING = "NOTIF_ACTIONABLE_WARNING",
  NOTIF_ACTIONABLE_ERROR = "NOTIF_ACTIONABLE_ERROR",
  NOTIF_ACTIONABLE_CRITICAL = "NOTIF_ACTIONABLE_CRITICAL",
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

export enum COMPONENT_TYPE {
  ALL = "ALL",
  ACCESS = "ACCESS",
  BACKHAUL = "BACKHAUL",
  POWER = "POWER",
  SWITCH = "SWITCH",
  SPECTRUM = "SPECTRUM",
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
      return NOTIFICATION_TYPE.NOTIF_INVALID;
    case 1:
      return NOTIFICATION_TYPE.NOTIF_INFO;
    case 2:
      return NOTIFICATION_TYPE.NOTIF_WARNING;
    case 3:
      return NOTIFICATION_TYPE.NOTIF_ERROR;
    case 4:
      return NOTIFICATION_TYPE.NOTIF_CRITICAL;
    default:
      return NOTIFICATION_TYPE.NOTIF_INVALID;
  }
};
