import { THeaders } from "../../common/types";
import AlertAPI from "../datasource/alert_api";

export interface Context {
  dataSources: {
    dataSource: AlertAPI;
  };
  headers: THeaders;
}
