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
export enum NODE_CONNECTIVITY {
  UNKNOWN = "unkown",
  OFFLINE = "offline",
  ONLINE = "online",
}
export enum NODE_TYPE {
  tnode = "tnode",
  anode = "anode",
  hnode = "hnode",
}
registerEnumType(NODE_STATUS, {
  name: "NodeStatusEnum",
  description: "Node status enums",
});
registerEnumType(NODE_CONNECTIVITY, {
  name: "NodeConnectivityEnum",
  description: "Node connectivity enums",
});
registerEnumType(NODE_TYPE, {
  name: "NodeTypeEnum",
  description: "Node type enums",
});
