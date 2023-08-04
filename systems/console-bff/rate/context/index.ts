import { THeaders } from "../../common/types";
import RateAPI from "../datasource/rate_api";

export interface Context {
  dataSources: {
    dataSource: RateAPI;
  };
  headers: THeaders;
}
