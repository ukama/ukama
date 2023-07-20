import NodeAPI from "../dataSource/node-api";

export interface Context {
  dataSources: {
    dataSource: NodeAPI;
  };
}
