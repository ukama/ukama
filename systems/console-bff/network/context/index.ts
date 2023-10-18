import { THeaders } from "../../common/types";
import NetworkAPI from "../datasource/network_api";

export interface Context {
  dataSources: {
    dataSource: NetworkAPI;
  };
  headers: THeaders;
}
