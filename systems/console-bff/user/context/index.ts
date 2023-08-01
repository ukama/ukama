import UserAPI from "../datasource/user_api";

export interface Context {
  dataSources: {
    dataSource: UserAPI;
  };
}
