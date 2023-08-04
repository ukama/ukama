import { THeaders } from "../../common/types";
import SubscriberAPI from "../datasource/subscriber_api";

export interface Context {
  dataSources: {
    dataSource: SubscriberAPI;
  };
  headers: THeaders;
}
