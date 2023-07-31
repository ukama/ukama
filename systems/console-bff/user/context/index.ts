import UserAPI from "../datasource/userapi";

export interface Context {
  dataSources: {
    dataSource: UserAPI;
  };
}
