import PackagesAPI from "../datasource/package_api";

export interface Context {
  dataSources: {
    dataSource: PackagesAPI;
  };
}
