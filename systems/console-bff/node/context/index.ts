import { THeaders } from "../../common/types";
import NodeAPI from "../dataSource/node-api";

export interface Context {
  dataSources: {
    dataSource: NodeAPI;
  };
  headers: THeaders;
}
