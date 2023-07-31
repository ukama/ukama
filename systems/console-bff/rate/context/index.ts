import RateAPI from "../datasource/rate_api";

export interface Context {
  dataSources: {
    dataSource: RateAPI;
  };
}
