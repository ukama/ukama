import { THeaders } from "../../common/types";
import SimAPI from "../datasource/sim_api";

export interface Context {
  dataSources: {
    dataSource: SimAPI;
  };
  headers: THeaders;
}
