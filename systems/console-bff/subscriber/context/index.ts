import SubscriberAPI from "../datasource/subscriber_api";

export interface Context {
  dataSources: {
    dataSource: SubscriberAPI;
  };
}
